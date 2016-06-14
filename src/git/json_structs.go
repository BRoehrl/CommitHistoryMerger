package git

import (
	"time"
)

// Repo is a struct for a single Repository of an Organization
type Repo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	HtmlURL        string `json:"html_url"`
	Language       string `json:"language"`
	DefaultBranch  string `json:"default_branch"`
	Branches       map[string]string
	SelectedBranch string
}

// Repos is a slice of Repo elements
type Repos []Repo

// User is a struct for a GitHub user
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"login"`
	URL     string `json:"url"`
	HtmlURL string `json:"html_url"`
}

// JSONVars is a struct for the Socket communication
type JSONVars struct {
	Author    string `json:"author"`
	Repo      string `json:"repo"`
	Querydate string `json:"date"`
}

// JSONCommit is a struct for a GitHub api commit response
type JSONCommit struct {
	Sha          string       `json:"sha"`
	ActualCommit SimpleCommit `json:"commit"`
	URL          string       `json:"url"`
	HtmlURL      string       `json:"html_url"`
	Author       User         `json:"author"`
	Committer    User         `json:"committer"`
}

// Commit is a struct for a single commit
type Commit struct {
	Sha,
	Repo,
	Branch,
	Author,
	CreatorLink,
	Link,
	Comment string
	Time time.Time
}

// Commits is a sortable slice of Commits
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

// SimpleUser is a struct containing the basic user information
type SimpleUser struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// SimpleCommit is a simplified struct for a commit
type SimpleCommit struct {
	Author    SimpleUser `json:"author"`
	Committer SimpleUser `json:"committer"`
	Message   string     `json:"message"`
	URL       string     `json:"url"`
}

// Branch is a struct  for a GitHub api branch response
type Branch struct {
	Name   string     `json:"name"`
	Commit JSONCommit `json:"commit"`
}

// UserCache is a struct containing all caches of a specific user
type UserCache struct {
	UserID        string
	CachedCommits Commits
	CachedShas    map[string]bool
	CachedAuthors map[string]bool
	CachedRepos   map[string]bool
	CacheTime     time.Time
	AllRepos      Repos
	UpdateCommits bool
	UpdateAll     bool
	Config        Config
}

// GetNewUserCache returns a new UserCache with defaultValues
func GetNewUserCache() (userCache UserCache) {
	userCache = UserCache{}
	userCache.CachedCommits = Commits{}
	userCache.CachedShas = make(map[string]bool)
	userCache.CachedAuthors = make(map[string]bool)
	userCache.CachedRepos = make(map[string]bool)
	userCache.CacheTime = time.Now().AddDate(0, 0, 1)
	userCache.AllRepos = Repos{}
	return
}

// SortableCommits is a sortable slice of JSONCommits
type SortableCommits []JSONCommit

func (sc SortableCommits) Len() int {
	return len(sc)
}
func (sc SortableCommits) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}
func (sc SortableCommits) Less(i, j int) bool {
	di := sc[i].ActualCommit.Author.Date
	dj := sc[j].ActualCommit.Author.Date
	return di.Before(dj)
}
