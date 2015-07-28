package processor

import (
	"git"
	"sort"
	"time"
)

type Commit struct {
	Sha,
	Repo,
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
}

func flushCommitCache() {
	cachedCommits = Commits{}
	cachedShas = make(map[string]bool)
	cachedAuthors = make(map[string]bool)
	cachedRepos = make(map[string]bool)
	cacheTime = time.Now()
}

func GetGitCommits(since time.Time) (err error) {
	allRepos, err := git.GetRepositories()
	if err != nil {
		return
	}
	for _, repo := range allRepos {
		singleRepoCommits, err := repo.GetAllCommitsUntil(since)
		if err != nil {
			return err
		}
		for _, gitCom := range singleRepoCommits {
			newCommit := Commit{
				Sha:     gitCom.Sha,
				Repo:    repo.Name,
				Author:  gitCom.ActualCommit.Author.Name,
				Link:    gitCom.HtmlUrl,
				Comment: gitCom.ActualCommit.Message,
				Time:    gitCom.ActualCommit.Author.Date,
			}
			addSingleCommitToCache(newCommit, false)
		}
	}
	sort.Sort(cachedCommits)
	if since.Before(cacheTime) {
		cacheTime = since
	}

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
