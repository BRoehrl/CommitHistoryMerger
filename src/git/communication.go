package git

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"fmt"
)

var (
	gitUrl,
	baseOrganisation,
	gitAuthkey,
	sinceTime string
	maxRepos,
	maxBranches,
	RateLimit,
	RateLimitRemaining,
	RateLimitReset int

	islastPage bool
)

func init() {
	gitUrl = "https://api.github.com"
	baseOrganisation = "informationgrid"
	gitAuthkey = ""
	maxRepos = 100
	maxBranches = 100

	twoMonthAgo := time.Now().AddDate(0, -2, 0)
	sinceTime = twoMonthAgo.Format("2006-01-02")
}

type Config struct {
	GitUrl, BaseOrganisation, GitAuthkey, SinceTime string
	MaxRepos, MaxBranches                           int
}

func SetConfig(connData Config) (settingsChanged bool) {
	if connData.GitUrl != "" && connData.GitUrl != gitUrl {
		gitUrl = connData.GitUrl
		settingsChanged = true
	}
	if connData.BaseOrganisation != "" && baseOrganisation != connData.BaseOrganisation {
		baseOrganisation = connData.BaseOrganisation
		settingsChanged = true
	}
	if strings.Replace(connData.GitAuthkey, "*", "", -1) != "" && gitAuthkey != connData.GitAuthkey {
		gitAuthkey = connData.GitAuthkey
		settingsChanged = true
	}
	if connData.SinceTime != "" && sinceTime != connData.SinceTime {
		sinceTime = connData.SinceTime
		settingsChanged = true
	}
	if connData.MaxRepos != 0 && maxRepos != connData.MaxRepos {
		maxRepos = connData.MaxRepos
		settingsChanged = true
	}
	if connData.MaxBranches != 0 && maxBranches != connData.MaxBranches {
		maxBranches = connData.MaxBranches
		settingsChanged = true
	}
	return settingsChanged
}

func GetConfig() Config {
	return Config{GitUrl: gitUrl, BaseOrganisation: baseOrganisation, GitAuthkey: "******", SinceTime: sinceTime, MaxRepos: maxRepos, MaxBranches: maxBranches}
}

func getResponse(url, baseAuthkey string) (resp *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Basic "+baseAuthkey)
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

func UnmarshalFromGetResponse(url string, i interface{}) (err error) {
	resp, err := getResponse(url, gitAuthkey)
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
	highestPageNumber := (maxRepos-1)/100 + 1
	fmt.Println("Repo:", highestPageNumber)
	for currentPage <=  highestPageNumber{
		repoQuery := gitUrl + "/orgs/" + baseOrganisation + "/repos?per_page=100&page=" + strconv.Itoa(currentPage)
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
	highestPageNumber := (maxBranches-1)/100 + 1
	branches := []Branch{}
	for currentPage <=  highestPageNumber{
		branchQuery := gitUrl + "/repos/" + baseOrganisation + "/" + repo.Name + "/branches?per_page=100&page=" + strconv.Itoa(currentPage)
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
	repo.Branches = branchMap
	return repo, err
}

// GetCommits returns the 100 newest commits for the specified repository
func (r Repo) GetCommits() (commits []JsonCommit, err error) {
	query := gitUrl + "/repos/" + baseOrganisation + "/" + r.Name + "/commits?per_page=100"
	err = UnmarshalFromGetResponse(query, &commits)
	return
}

// GetFirstNCommits returns the N newest commits for the specified repository.
// Note that n/100 querries are sent to the server
func (r Repo) GetFirstNCommits(n int) (commits []JsonCommit, err error) {
	currentPage := 1
	for {
		query := gitUrl + "/repos/" + baseOrganisation + "/" + r.Name + "/commits?per_page=100&page=" + strconv.Itoa(currentPage)
		currentPage++
		var singlePage []JsonCommit
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

// GetAllCommitsUntil returns all commits commited before Date to and after Date from
// Note that for each 100 querries a new request is sent
func (r Repo) GetAllCommitsBetween(from, to time.Time) (commits []JsonCommit, err error) {
	currentPage := 1
	for {
		query := gitUrl + "/repos/" + baseOrganisation
		query += "/" + r.Name
		query += "/commits?since=" + from.Format(time.RFC3339) + "&until=" + to.Format(time.RFC3339)
		query += "&sha=" + r.Branches[r.SelectedBranch]
		query += "&per_page=100&page=" + strconv.Itoa(currentPage)

		currentPage++
		var singlePage []JsonCommit
		err = UnmarshalFromGetResponse(query, &singlePage)
		commits = append(commits, singlePage...)

		if err != nil || islastPage {
			break
		}

	}
	return
}

// GetAllCommitsUntil returns all commits commited after Date d
// Note that for each 100 querries a new request is sent
func (r Repo) GetAllCommitsUntil(d time.Time) (commits []JsonCommit, err error) {
	return r.GetAllCommitsBetween(time.Now(), d)
}
