package processor

import (
	"errors"
	"git"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
	"user"
)

// A Query to the Backend
type Query struct {
	Commit   string
	Authors  []string
	Repos    []string
	Since    time.Time
	UseRegex bool
}

// SendCommits sends all commits matching the query to a supplied channel
func SendCommits(userID string, query Query, commits chan git.Commit) {
	uc := user.GetUserCache(userID)

	if uc.UpdateCommits || uc.UpdateAll {
		flushCommitCache(&uc)
		uc.UpdateCommits = false
		if uc.UpdateAll {
			flushRepos(&uc)
		}
		uc.UpdateAll = false
	}
	// if no date set use default Date
	if query.Since.Equal(time.Time{}) {
		query.Since = uc.Config.SinceTime
	}

	allCommits := make(chan git.Commit)
	// fetch commits if not in cache else send cache to channel
	if query.Since.Before(uc.CacheTime) {
		go sendGitCommits(&uc, query.Since, uc.CacheTime, allCommits)
		uc.CacheTime = query.Since
	} else {
		allCommits = make(chan git.Commit, len(uc.CachedCommits))
		for _, commit := range uc.CachedCommits {
			allCommits <- commit
		}
		close(allCommits)
	}

	for commit := range allCommits {
		addSingleCommitToCache(&uc, commit, false)
		if keepCommit(query, commit) {
			commits <- commit
		}
	}
	user.SetCachedCommits(userID, uc.CachedCommits)
	close(commits)
	return
}

func keepCommit(query Query, commit git.Commit) bool {

	keep := true
	if commit.Time.Before(query.Since) {
		keep = false
	}

	if keep {
		if len(query.Authors) != 0 {
			keep = false
			for _, author := range query.Authors {
				if commit.Author == author || author == "" {
					keep = true
					break
				}
				if query.UseRegex {
					matched, _ := regexp.MatchString(strings.ToLower(author), strings.ToLower(commit.Author))
					if matched {
						keep = true
						break
					}
				}
			}
		}
	}

	if keep {
		if len(query.Repos) != 0 {
			keep = false
			for _, repo := range query.Repos {
				if commit.Repo == repo || repo == "" {
					keep = true
					break
				}
				if query.UseRegex {
					matched, _ := regexp.MatchString(strings.ToLower(repo), strings.ToLower(commit.Repo))
					if matched {
						keep = true
						break
					}
				}
			}
		}
	}

	if keep {
		keep = false
		if commit.Comment == query.Commit || query.Commit == "" {
			keep = true
		} else if query.UseRegex {
			keep, _ = regexp.MatchString(strings.ToLower(query.Commit), strings.ToLower(commit.Comment))
		}
	}

	return keep
}

// GetSingleCommit returns the commit with sha if in cache
func GetSingleCommit(userCache git.UserCache, sha string) (singleCommit git.Commit) {
	if !userCache.CachedShas[sha] {
		return
	}
	for _, com := range userCache.CachedCommits {
		if com.Sha == sha {
			singleCommit = com
			return
		}
	}
	return
}

// GetCacheTimeString returns the earliest date for which the commits are cached as a string
func GetCacheTimeString(userCache *git.UserCache) (cacheTimeString string) {
	return userCache.CacheTime.Format(time.RFC3339)[:10]
}

// GetCachedAuthors returns all cached authornames
func GetCachedAuthors(userCache *git.UserCache) (authors []string) {
	for key := range userCache.CachedAuthors {
		authors = append(authors, key)
	}
	sort.Strings(authors)
	return
}

// GetCachedRepos returns all cached repositorynames
func GetCachedRepos(userCache *git.UserCache) (repos []string) {
	for key := range userCache.CachedRepos {
		repos = append(repos, key)
	}
	sort.Strings(repos)
	return
}

// GetCachedRepoObjects returns all cached repositories
func GetCachedRepoObjects(userCache *git.UserCache) (repos git.Repos, err error) {
	if len(userCache.CachedRepos) == 0 {
		userCache.AllRepos, err = git.GetRepositories(userCache.Config)
	}
	for _, repo := range userCache.AllRepos {
		userCache.CachedRepos[repo.Name] = true
	}
	return userCache.AllRepos, err
}

// SetRepoBranch sets the branch of a repository
func SetRepoBranch(userCache git.UserCache, repoName, branchName string) (err error) {
	if !userCache.CachedRepos[repoName] {
		return errors.New("Repository not found/cached: " + repoName)
	}
	for i, repo := range userCache.AllRepos {
		if repo.Name == repoName {
			for branch := range repo.Branches {
				if branch == branchName {
					if repo.SelectedBranch != branch {
						repo.SelectedBranch = branch
						userCache.UpdateCommits = true
						userCache.AllRepos[i] = repo
					}
					return
				}
			}
			return errors.New("Repository found, but no Branch named: " + branchName)
		}
	}
	return errors.New("Repository found, but not in allRepos: " + repoName)
}

// UpdateDefaultBranch sets the Branch of all repositories to the default branch from the config.
func UpdateDefaultBranch(userCache git.UserCache) {
	defaultBranch := userCache.Config.MiscDefaultBranch
	for i, repo := range userCache.AllRepos {
		if repo.Branches[defaultBranch] != "" {
			repo.SelectedBranch = defaultBranch
			userCache.AllRepos[i] = repo
		}
	}
}

// SetConfig sets the config and updates the default branch if necessary
func SetConfig(userCache git.UserCache, config git.Config) {
	completeUpdate, miscBranchChanged := git.SetConfig(userCache.Config, config)
	userCache.UpdateAll = completeUpdate
	if miscBranchChanged {
		UpdateDefaultBranch(userCache)
		userCache.UpdateCommits = true
	}
}

// GetSavedConfigs returns all saved configfilenames
func GetSavedConfigs() (fileNames []string) {
	fileNames, err := getSavedConfigs()
	if err != nil {
		log.Println(err)
	}
	return
}
