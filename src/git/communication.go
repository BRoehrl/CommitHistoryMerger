package git

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ResponseMetaData contains all necessary Metadata from a GitHub API response
type ResponseMetaData struct {
	RateLimit,
	RateLimitRemaining,
	RateLimitReset int
	IslastPage bool
}

// AuthTokenToLastResponse maps the metadata of the last response to the respective AuthToken
var AuthTokenToLastResponse map[string](ResponseMetaData)

func init() {
	AuthTokenToLastResponse = make(map[string](ResponseMetaData))
}

// GetAuthKeyFromGit TODO
func GetAuthKeyFromGit(code, clientID, clientSecret string) (string, error) {
	client := &http.Client{}
	form := url.Values{}
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	form.Add("code", code)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(req, err)
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	type AuthtokenResponse struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	auRe := AuthtokenResponse{}
	err = json.Unmarshal([]byte(body), &auRe)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return auRe.AccessToken, nil
}

// GetUserFromToken returns the authTokens owner
func GetUserFromToken(authToken string) User {
	CurrentUser := User{}
	UnmarshalFromGetResponse(DefaultConfig.GitURL+"/user", authToken, &CurrentUser)
	return CurrentUser
}

func getResponse(url, baseAuthkey string) (resp *http.Response, responseMeta ResponseMetaData, err error) {
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
	responseMeta.RateLimit, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	responseMeta.RateLimitRemaining, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	responseMeta.RateLimitReset, err = strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))
	responseMeta.IslastPage = true
	//check if only one page
	if link := resp.Header.Get("Link"); link != "" {
		//check if on last page
		if strings.Contains(link, "rel=\"next\"") {
			responseMeta.IslastPage = false
		}
	}
	AuthTokenToLastResponse[baseAuthkey] = responseMeta
	return
}

// UnmarshalFromGetResponse unmarshals the json response of a git api call
// into an interface i
func UnmarshalFromGetResponse(url, authKey string, i interface{}) (responseMetaData ResponseMetaData, err error) {
	resp, responseMetaData, err := getResponse(url, authKey)
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

// GetRepositories returns all or userConfig.MaxRepos repositories of the baseOrganisation.
func GetRepositories(userConfig Config) (allRepos Repos, err error) {
	var responseMetaData ResponseMetaData
	currentPage := 1
	highestPageNumber := (userConfig.MaxRepos-1)/100 + 1
	for currentPage <= highestPageNumber {
		repoQuery := userConfig.GitURL + "/orgs/" + userConfig.BaseOrganisation + "/repos?per_page=100&page=" + strconv.Itoa(currentPage)
		currentPage++
		var reposPage Repos
		responseMetaData, err = UnmarshalFromGetResponse(repoQuery, userConfig.GitAuthkey, &reposPage)
		if err != nil {
			return
		}
		allRepos = append(allRepos, reposPage...)

		if responseMetaData.IslastPage {
			break
		}
	}
	allRepos, err = AddBranchesToRepos(allRepos, userConfig)
	return
}

var branchWaitGroup sync.WaitGroup

// AddBranchesToRepos adds to each Repository its branches
func AddBranchesToRepos(allRepos Repos, userConfig Config) (reposWithBranches Repos, err error) {
	repoChannel := make(chan Repo, len(allRepos))
	for _, repo := range allRepos {
		go addBranchesToSingleRepo(userConfig, repo, repoChannel)
	}

	for repo := range repoChannel {
		if len(reposWithBranches) >= len(allRepos)-1 {
			break
		}
		reposWithBranches = append(reposWithBranches, repo)
	}
	return
}

func addBranchesToSingleRepo(userConfig Config, repo Repo, repoChannel chan Repo) (r Repo, err error) {
	var responseMetaData ResponseMetaData
	currentPage := 1
	highestPageNumber := (userConfig.MaxBranches-1)/100 + 1
	branches := []Branch{}
	for currentPage <= highestPageNumber {
		branchQuery := userConfig.GitURL + "/repos/" + userConfig.BaseOrganisation + "/" + repo.Name + "/branches?per_page=100&page=" + strconv.Itoa(currentPage)
		currentPage++
		var branchesPage []Branch
		responseMetaData, err = UnmarshalFromGetResponse(branchQuery, userConfig.GitAuthkey, &branchesPage)
		branches = append(branches, branchesPage...)
		if err != nil {
			log.Println(err)
			return repo, err
		}
		if responseMetaData.IslastPage {
			break
		}
	}
	branchMap := make(map[string]string)
	for _, b := range branches {
		branchMap[b.Name] = b.Commit.Sha
	}
	repo.SelectedBranch = repo.DefaultBranch
	if branchMap[userConfig.MiscDefaultBranch] != "" {
		repo.SelectedBranch = userConfig.MiscDefaultBranch
	}
	repo.Branches = branchMap
	repoChannel <- repo
	if err != nil {
		log.Println(err)
	}
	return repo, err
}

// CommitWaitGroup is a WaitGroup for sending commits via SendAllCommitsBetween
var CommitWaitGroup sync.WaitGroup

// SendAllCommitsBetween sends all commits of repo r between times from and to the the supplied channel
func (r Repo) SendAllCommitsBetween(from, to time.Time, allComits chan Commit, userConfig Config) {
	currentPage := 1
	for {
		query := userConfig.GitURL + "/repos/" + userConfig.BaseOrganisation
		query += "/" + r.Name
		query += "/commits?since=" + from.Format(time.RFC3339) + "&until=" + to.Format(time.RFC3339)
		query += "&sha=" + r.Branches[r.SelectedBranch]
		query += "&per_page=100&page=" + strconv.Itoa(currentPage)

		currentPage++
		var singlePage []JSONCommit
		responseMetaData, err := UnmarshalFromGetResponse(query, userConfig.GitAuthkey, &singlePage)
		for _, gitCom := range singlePage {
			newCommit := Commit{
				Sha:         gitCom.Sha,
				Repo:        r.Name,
				Branch:      r.SelectedBranch,
				Author:      gitCom.ActualCommit.Author.Name,
				CreatorLink: gitCom.Author.HtmlURL,
				Link:        gitCom.HtmlURL,
				Comment:     gitCom.ActualCommit.Message,
				Time:        gitCom.ActualCommit.Author.Date,
			}
			allComits <- newCommit
		}

		if err != nil || responseMetaData.IslastPage {
			break
		}
	}
	CommitWaitGroup.Done()
	return
}
