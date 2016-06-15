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
	"user"
)

var defaultUserCache git.UserCache

// LoadedConfig is the name of the loaded config
var LoadedConfig string

func init() {
	defaultUserCache = git.GetNewUserCache()
	defaultUserCache.CachedRepos = make(map[string]bool)
}

func flushCommitCache(userCache *git.UserCache) {
	userCache.CachedCommits = git.Commits{}
	userCache.CachedShas = make(map[string]bool)
	userCache.CachedAuthors = make(map[string]bool)
	userCache.CacheTime = time.Now().AddDate(0, 0, 1)
}

func flushRepos(userCache *git.UserCache) {
	userCache.AllRepos = git.Repos{}
	userCache.CachedRepos = make(map[string]bool)
}

func sendGitCommits(userCache *git.UserCache, from, to time.Time, allCommits chan git.Commit) {
	if len(userCache.AllRepos) == 0 {
		userCache.AllRepos, _ = git.GetRepositories(userCache.Config)
	}
	for _, repo := range userCache.AllRepos {
		git.CommitWaitGroup.Add(1)
		go repo.SendAllCommitsBetween(from, to, allCommits, userCache.Config)
	}
	git.CommitWaitGroup.Wait()
	close(allCommits)
	sort.Sort(userCache.CachedCommits)
	user.SetUserCache(userCache.UserID, *userCache)
	return
}

func addSingleCommitToCache(userCache *git.UserCache, nc git.Commit, reSort bool) (*git.UserCache, bool) {
	commitAdded := !userCache.CachedShas[nc.Sha]
	if commitAdded {
		userCache.CachedShas[nc.Sha] = true
		userCache.CachedAuthors[nc.Author] = true
		userCache.CachedRepos[nc.Repo] = true
		userCache.CachedCommits = append(userCache.CachedCommits, nc)
	}
	if reSort {
		sort.Sort(userCache.CachedCommits)
	}
	return userCache, commitAdded
}

type completeConfig struct {
	Baseconfig       git.Config
	SelectedBranches map[string]string
}

// SaveCompleteConfig saves the current config as a file
func SaveCompleteConfig(userCache git.UserCache, fileName string) error {
	baseConfig := userCache.Config
	selectedBranches := make(map[string]string)
	for _, repo := range userCache.AllRepos {
		selectedBranches[repo.Name] = repo.SelectedBranch
	}
	err := saveInJSONFile(completeConfig{baseConfig, selectedBranches}, "configs", fileName)
	if err != nil {
		return err
	}
	LoadedConfig = fileName
	return nil
}

func getSavedConfigs() (fileNames []string, err error) {
	file, err := os.Open("configs/")
	if err != nil {
		log.Println("Config-folder not found", "configs/")
		return
	}
	return file.Readdirnames(0)
}

// LoadCompleteConfig loads a config-file
func LoadCompleteConfig(userID string) (err error) {
	file, err := os.Open("configs/" + userID)
	if err != nil {
		log.Println("Config-file not found", "configs/"+userID)
		return
	}
	decoder := json.NewDecoder(file)
	completeConfig := completeConfig{}
	err = decoder.Decode(&completeConfig)
	if err != nil {
		log.Println("Could not parse Config-file", file.Name())
	}
	user.AddUser(userID, completeConfig.Baseconfig.GitAuthkey)
	SetConfig(userID, completeConfig.Baseconfig)
	//reload repositories
	GetCachedRepoObjects(userID)
	for repo, branch := range completeConfig.SelectedBranches {
		SetRepoBranch(userID, repo, branch)
	}
	uc := user.GetUserCache(userID)
	flushCommitCache(&uc)
	return
}

func saveInJSONFile(i interface{}, dir string, fileName string) (err error) {
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
