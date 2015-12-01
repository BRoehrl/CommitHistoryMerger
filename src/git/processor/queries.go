package processor

import (
	"errors"
	"fmt"
	"git"
	"log"
	"sort"
	"time"
)

type Query struct {
	Authors []string
	Repos   []string
	Since   time.Time
}

var updateCommits, updateAll bool

func SendCommits(query Query, commits chan git.Commit) {

	if updateCommits || updateAll {
		flushCommitCache()
		updateCommits = false
		if updateAll {
			flushRepos()
		}
		updateAll = false
	}
	// if no date set use default Date
	if query.Since.Equal(time.Time{}) {
		query.Since = git.GetConfig().SinceTime
	}

	allCommits := make(chan git.Commit)
	// fetch commits if not in cache else send cache to channel
	if query.Since.Before(cacheTime) {
		go sendGitCommits(query.Since, cacheTime, allCommits)
		cacheTime = query.Since
	} else {
		allCommits = make(chan git.Commit, len(cachedCommits))
		for _, commit := range cachedCommits {
			allCommits <- commit
		}
		close(allCommits)
	}
	for commit := range allCommits {
		addSingleCommitToCache(commit, false)
		if keepCommit(query, commit) {
			commits <- commit
		}
	}
	close(commits)
	return
}

func keepCommit(query Query, commit git.Commit) bool {
	keep := true
	if len(query.Authors) != 0 {
		keep = false
		for _, author := range query.Authors {
			if commit.Author == author {
				keep = true
			}
		}
	}

	if keep {
		if len(query.Repos) != 0 {
			keep = false
			for _, repo := range query.Repos {
				if commit.Repo == repo {
					keep = true
				}
			}
		}
	}
	if keep {
		if commit.Time.Before(query.Since) {
			keep = false
		}
	}
	return keep
}

func GetCommits(query Query) (commits Commits) {
	commits = Commits{}
	if updateCommits || updateAll {
		flushCommitCache()
		updateCommits = false
		if updateAll {
			flushRepos()
		}
		updateAll = false
	}
	// if no date set use default Date
	if query.Since.Equal(time.Time{}) {
		query.Since = git.GetConfig().SinceTime
	}
	if query.Since.Before(cacheTime) {
		err := getGitCommits(query.Since, cacheTime)
		if err != nil {
			fmt.Println(err)
		}
		cacheTime = query.Since
	}
	for _, commit := range cachedCommits {
		keep := true
		if len(query.Authors) != 0 {
			keep = false
			for _, author := range query.Authors {
				if commit.Author == author {
					keep = true
				}
			}
		}

		if keep {
			if len(query.Repos) != 0 {
				keep = false
				for _, repo := range query.Repos {
					if commit.Repo == repo {
						keep = true
					}
				}
			}
		}

		if keep {
			if commit.Time.Before(query.Since) {
				keep = false
			}
		}

		if keep {
			commits = append(commits, commit)
		}
	}
	return
}

func GetSingleCommit(sha string) (singleCommit git.Commit) {
	if !cachedShas[sha] {
		return
	}
	for _, com := range cachedCommits {
		if com.Sha == sha {
			singleCommit = com
			return
		}
	}
	return
}

func GetCachedAuthors() (authors []string) {
	for key := range cachedAuthors {
		authors = append(authors, key)
	}
	sort.Strings(authors)
	return
}
func GetCachedRepos() (repos []string) {
	for key := range cachedRepos {
		repos = append(repos, key)
	}
	sort.Strings(repos)
	return
}
func GetCachedRepoObjects() (repos git.Repos, err error) {
	if len(cachedRepos) == 0 {
		allRepos, err = git.GetRepositories()
	}
	for _, repo := range allRepos {
		cachedRepos[repo.Name] = true
	}
	return allRepos, err
}

func SetRepoBranch(repoName, branchName string) (err error) {
	if !cachedRepos[repoName] {
		return errors.New("Repository not found/cached: " + repoName)
	}
	for i, repo := range allRepos {
		if repo.Name == repoName {
			for branch := range repo.Branches {
				if branch == branchName {
					if repo.SelectedBranch != branch {
						repo.SelectedBranch = branch
						updateCommits = true
						allRepos[i] = repo
					}
					return
				}
			}
			return errors.New("Repository found, but no Branch named: " + branchName)
		}
	}
	return errors.New("Repository found, but not in allRepos: " + repoName)
}

func UpdateDefaultBranch() {
	defaultBranch := git.GetConfig().MiscDefaultBranch
	for i, repo := range allRepos {
		if repo.Branches[defaultBranch] != "" {
			repo.SelectedBranch = defaultBranch
			allRepos[i] = repo
		}
	}
}

func SetConfig(config git.Config) {
	completeUpdate, miscBranchChanged := git.SetConfig(config)
	updateAll = completeUpdate
	if miscBranchChanged {
		UpdateDefaultBranch()
		updateCommits = true
	}
}

func GetSavedConfigs() (fileNames []string) {
	fileNames, err := getSavedConfigs()
	if err != nil {
		log.Println(err)
	}
	return
}
