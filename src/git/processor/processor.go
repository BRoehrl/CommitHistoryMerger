package processor

import (
	"encoding/json"
	"fmt"
	"git"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

type Commit struct {
	Sha,
	Repo,
	Branch,
	Author,
	Link,
	Comment string
	Time time.Time
}

var cachedCommits Commits
var cachedShas map[string]bool
var cachedAuthors map[string]bool
var cachedRepos map[string]bool
var cacheTime time.Time
var allRepos git.Repos

type Commits []Commit

func (c Commits) Len() int {
	return len(c)
}
func (c Commits) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Commits) Less(i, j int) bool {
	return c[i].Time.After(c[j].Time)
}

func init() {
	flushCommitCache()
	cachedRepos = make(map[string]bool)
}

func flushCommitCache() {
	cachedCommits = Commits{}
	cachedShas = make(map[string]bool)
	cachedAuthors = make(map[string]bool)
	cacheTime = time.Now()
}

func flushRepos() {
	allRepos = git.Repos{}
	cachedRepos = make(map[string]bool)
}

func getGitCommits(from, to time.Time) (err error) {
	if len(cachedRepos) == 0 {
		allRepos, err = git.GetRepositories()
	}
	if err != nil {
		return
	}
	for _, repo := range allRepos {
		singleRepoCommits, err := repo.GetAllCommitsBetween(from, to)
		if err != nil {
			return err
		}
		for _, gitCom := range singleRepoCommits {
			newCommit := Commit{
				Sha:     gitCom.Sha,
				Repo:    repo.Name,
				Branch:  repo.SelectedBranch,
				Author:  gitCom.ActualCommit.Author.Name,
				Link:    gitCom.HtmlUrl,
				Comment: gitCom.ActualCommit.Message,
				Time:    gitCom.ActualCommit.Author.Date,
			}
			addSingleCommitToCache(newCommit, false)
		}
	}
	sort.Sort(cachedCommits)
	return
}

func addCommitsToCache(newCommits Commits) {
	for _, nc := range newCommits {
		addSingleCommitToCache(nc, false)
	}
	sort.Sort(cachedCommits)
}

func addSingleCommitToCache(nc Commit, reSort bool) (commitAdded bool) {
	commitAdded = !cachedShas[nc.Sha]
	if commitAdded {
		cachedShas[nc.Sha] = true
		cachedAuthors[nc.Author] = true
		cachedRepos[nc.Repo] = true
		cachedCommits = append(cachedCommits, nc)

	}
	if reSort {
		sort.Sort(cachedCommits)
	}
	return
}

type completeConfig struct {
	Baseconfig       git.Config
	SelectedBranches map[string]string
}

func saveCompleteConfig(fileName string) {
	baseConfig := git.GetConfig()
	selectedBranches := make(map[string]string)
	for _, repo := range allRepos {
		selectedBranches[repo.Name] = repo.SelectedBranch
	}
	saveInJsonFile(completeConfig{baseConfig, selectedBranches}, "configs", fileName)
}

func getSavedConfigs() (fileNames []string, err error) {
	file, err := os.Open("configs/")
	if err != nil {
		log.Println("Config-folder not found", "configs/")
		return
	}
	return file.Readdirnames(0)
}

func loadCompleteConfig(fileName string) (err error) {
	file, err := os.Open("configs/" + fileName)
	if err != nil {
		log.Println("Config-file not found", "configs/"+fileName)
		return
	}
	decoder := json.NewDecoder(file)
	completeConfig := completeConfig{}
	err = decoder.Decode(&completeConfig)
	if err != nil {
		log.Println("Could not parse Config-file", file.Name())
	}
	log.Println("complete config:", completeConfig)

	SetConfig(completeConfig.Baseconfig)
	//reload repositories
	GetCachedRepoObjects()
	for repo, branch := range completeConfig.SelectedBranches {
		SetRepoBranch(repo, branch)
	}
	flushCommitCache()
	return
}

func saveInJsonFile(i interface{}, dir string, fileName string) (err error) {
	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		} else {
			log.Println(err)
			return
		}
	}

	path := fmt.Sprint(dir, "/", fileName)
	os.Remove(path)

	b, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
		return
	}

	ioutil.WriteFile(path, b, 0644)
	return
}
