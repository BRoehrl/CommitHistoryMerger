package processor

import (
	"errors"
	"fmt"
	"git"
	"sort"
	"time"
)

type Query struct {
	Authors []string
	Repos   []string
	Since   time.Time
}

var aBranchChanged, settingsChanged bool

func GetCommits(query Query) (commits Commits) {
	commits = Commits{}
	if aBranchChanged || settingsChanged {
		flushCommitCache()
		settingsChanged = false
		aBranchChanged = false
	}
	if query.Since.Before(cacheTime) {
		err := GetGitCommits(query.Since, cacheTime)
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
						aBranchChanged = true
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

func SetConfig(config git.Config) {
	if  git.SetConfig(config) {
		settingsChanged = true
	}
	fmt.Println(settingsChanged)
	fmt.Println(config)
}
