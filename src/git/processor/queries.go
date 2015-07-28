package processor

import (
	"fmt"
	"time"
	"sort"
)

type Query struct {
	Authors []string
	Repos   []string
	Since   time.Time
}

func GetCommits(query Query) (commits Commits) {
	commits = Commits{}
	if query.Since.Before(cacheTime) {
		err := GetGitCommits(query.Since)
		if err != nil {
			fmt.Println(err)
		}
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
