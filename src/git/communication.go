package git

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RateLimit is the RateLimit for GitHub API queries
var RateLimit int

// RateLimitRemaining is the number of request remaining until the next ratelimit reset
var RateLimitRemaining int

// RateLimitReset is the time left until the next ratelimit reset
var RateLimitReset int
var islastPage bool

func getResponse(url, baseAuthkey string) (resp *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "token "+baseAuthkey)
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	RateLimit, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	RateLimitRemaining, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	RateLimitReset, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))
	RateLimitRemaining, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	RateLimitReset, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))
	islastPage = true
	//check if only one page
	if link := resp.Header.Get("Link"); link != "" {
		//check if on last page
		if strings.Contains(link, "rel=\"next\"") {
			islastPage = false
		}
	}

	return
}

// UnmarshalFromGetResponse unmarshals the json response of a git api call
// into an interface i
func UnmarshalFromGetResponse(url string, i interface{}) (err error) {
	resp, err := getResponse(url, config.GitAuthkey)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(body), &i)
	return
}

// GetRepositories returns all or the first 100 repositories of the baseOrganisation.
func GetRepositories() (allRepos Repos, err error) {
	currentPage := 1
	highestPageNumber := (config.MaxRepos-1)/100 + 1
	for currentPage <= highestPageNumber {
		repoQuery := config.GitURL + "/orgs/" + config.BaseOrganisation + "/repos?per_page=100&page=" + strconv.Itoa(currentPage)
		currentPage++
		var reposPage Repos
		err = UnmarshalFromGetResponse(repoQuery, &reposPage)
		if err != nil {
			return
		}
		allRepos = append(allRepos, reposPage...)

		if islastPage {
			break
		}
	}
	allRepos, err = AddBranchesToRepos(allRepos)
	return
}

// AddBranchesToRepos adds to each Repository its branches
func AddBranchesToRepos(allRepos Repos) (reposWithBranches Repos, err error) {
	for _, repo := range allRepos {
		repo, err = addBranchesToSingleRepo(repo)
		if err != nil {
			return
		}
		reposWithBranches = append(reposWithBranches, repo)
	}
	return
}

func addBranchesToSingleRepo(repo Repo) (r Repo, err error) {
	currentPage := 1
	highestPageNumber := (config.MaxBranches-1)/100 + 1
	branches := []Branch{}
	for currentPage <= highestPageNumber {
		branchQuery := config.GitURL + "/repos/" + config.BaseOrganisation + "/" + repo.Name + "/branches?per_page=100&page=" + strconv.Itoa(currentPage)
		currentPage++
		var branchesPage []Branch
		err = UnmarshalFromGetResponse(branchQuery, &branchesPage)
		branches = append(branches, branchesPage...)
		if err != nil {
			return repo, err
		}
		if islastPage {
			break
		}
	}
	branchMap := make(map[string]string)
	for _, b := range branches {
		branchMap[b.Name] = b.Commit.Sha
	}
	repo.SelectedBranch = repo.DefaultBranch
	if branchMap[config.MiscDefaultBranch] != "" {
		repo.SelectedBranch = config.MiscDefaultBranch
	}
	repo.Branches = branchMap
	return repo, err
}

// GetCommits returns the 100 newest commits for the specified repository
func (r Repo) GetCommits() (commits []JSONCommit, err error) {
	query := config.GitURL + "/repos/" + config.BaseOrganisation + "/" + r.Name + "/commits?per_page=100"
	err = UnmarshalFromGetResponse(query, &commits)
	return
}

// GetFirstNCommits returns the N newest commits for the specified repository.
// Note that n/100 queries are sent to the server
func (r Repo) GetFirstNCommits(n int) (commits []JSONCommit, err error) {
	currentPage := 1
	for {
		query := config.GitURL + "/repos/" + config.BaseOrganisation + "/" + r.Name + "/commits?per_page=100&page=" + strconv.Itoa(currentPage)
		currentPage++
		var singlePage []JSONCommit
		err = UnmarshalFromGetResponse(query, &singlePage)
		commits = append(commits, singlePage...)

		if len(commits) >= n {
			commits = commits[:n]
			break
		}
		if err != nil || islastPage {
			break
		}

	}
	return
}

// GetAllCommitsBetween returns all commits commited before Date to and after Date from
// Note that for each 100 queries a new request is sent
func (r Repo) GetAllCommitsBetween(from, to time.Time) (commits []JSONCommit, err error) {
	currentPage := 1
	for {
		query := config.GitURL + "/repos/" + config.BaseOrganisation
		query += "/" + r.Name
		query += "/commits?since=" + from.Format(time.RFC3339) + "&until=" + to.Format(time.RFC3339)
		query += "&sha=" + r.Branches[r.SelectedBranch]
		query += "&per_page=100&page=" + strconv.Itoa(currentPage)

		currentPage++
		var singlePage []JSONCommit
		err = UnmarshalFromGetResponse(query, &singlePage)
		commits = append(commits, singlePage...)

		if err != nil || islastPage {
			break
		}

	}
	return
}

// CommitWaitGroup is a WaitGroup for sending commits via SendAllCommitsBetween
var CommitWaitGroup sync.WaitGroup

// SendAllCommitsBetween sends all commits of repo r between times from and to the the supplied channel
func (r Repo) SendAllCommitsBetween(from, to time.Time, allComits chan Commit) {
	currentPage := 1
	for {
		query := config.GitURL + "/repos/" + config.BaseOrganisation
		query += "/" + r.Name
		query += "/commits?since=" + from.Format(time.RFC3339) + "&until=" + to.Format(time.RFC3339)
		query += "&sha=" + r.Branches[r.SelectedBranch]
		query += "&per_page=100&page=" + strconv.Itoa(currentPage)

		currentPage++
		var singlePage []JSONCommit
		err := UnmarshalFromGetResponse(query, &singlePage)
		for _, gitCom := range singlePage {
			newCommit := Commit{
				Sha:     gitCom.Sha,
				Repo:    r.Name,
				Branch:  r.SelectedBranch,
				Author:  gitCom.ActualCommit.Author.Name,
				Link:    gitCom.HtmlURL,
				Comment: gitCom.ActualCommit.Message,
				Time:    gitCom.ActualCommit.Author.Date,
			}
			allComits <- newCommit
		}

		if err != nil || islastPage {
			break
		}
	}
	CommitWaitGroup.Done()
	return
}

// GetAllCommitsUntil returns all commits commited after Date d
// Note that for each 100 queries a new request is sent
func (r Repo) GetAllCommitsUntil(d time.Time) (commits []JSONCommit, err error) {
	return r.GetAllCommitsBetween(time.Now(), d)
}
