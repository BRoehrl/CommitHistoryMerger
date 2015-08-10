package git

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	gitUrl             string
	baseOrganisation   string
	gitAuthkey         string
	RateLimit          int
	RateLimitRemaining int
	RateLimitReset     int

	islastPage bool
)

func init() {
	gitUrl = "https://api.github.com"
	baseOrganisation = "/informationgrid"
	gitAuthkey = ""
}

type Config struct {
	GitUrl, BaseOrganisation, GitAuthkey string
}

func SetConfig(connData Config) {
	if connData.GitUrl != "" {
		gitUrl = connData.GitUrl
	}
	if connData.BaseOrganisation != "" {
		baseOrganisation = connData.BaseOrganisation
	}
	if connData.GitAuthkey != "" {
		gitAuthkey = connData.GitAuthkey
	}
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
	repoQuery := gitUrl + "/orgs" + baseOrganisation + "/repos?per_page=100"
	err = UnmarshalFromGetResponse(repoQuery, &allRepos)
	if err != nil {
		return
	}
	allRepos, err = AddBranchesToRepos(allRepos)
	return
}

func AddBranchesToRepos(allRepos Repos) (reposWithBranches Repos, err error) {
	for _, repo := range allRepos {
		branchQuery := gitUrl + "/repos" + baseOrganisation + "/" + repo.Name + "/branches?per_page=100"
		branches := []Branch{}
		err = UnmarshalFromGetResponse(branchQuery, &branches)
		if err != nil {
			return
		}
		branchMap := make(map[string]string)
		for _, b := range branches {
			branchMap[b.Name] = b.Commit.Sha
		}
		repo.SelectedBranch = repo.DefaultBranch
		repo.Branches = branchMap
		reposWithBranches = append(reposWithBranches, repo)
	}
	return
}

// GetCommits returns the 100 newest commits for the specified repository
func (r Repo) GetCommits() (commits []JsonCommit, err error) {
	query := gitUrl + "/repos" + baseOrganisation + "/" + r.Name + "/commits?per_page=100"
	err = UnmarshalFromGetResponse(query, &commits)
	return
}

// GetFirstNCommits returns the N newest commits for the specified repository.
// Note that n/100 querries are sent to the server
func (r Repo) GetFirstNCommits(n int) (commits []JsonCommit, err error) {
	currentPage := 1
	for {
		query := gitUrl + "/repos" + baseOrganisation + "/" + r.Name + "/commits?per_page=100&page=" + strconv.Itoa(currentPage)
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
		query := gitUrl + "/repos" + baseOrganisation
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
