package git

import (
	"time"
)

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

type Repos []Repo

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"login"`
	URL     string `json:"url"`
	HtmlURL string `json:"html_url"`
}

type JSONCommit struct {
	Sha          string       `json:"sha"`
	ActualCommit SimpleCommit `json:"commit"`
	URL          string       `json:"url"`
	HtmlURL      string       `json:"html_url"`
	Author       User         `json:"author"`
	Committer    User         `json:"committer"`
}

type Commit struct {
	Sha,
	Repo,
	Branch,
	Author,
	Link,
	Comment string
	Time time.Time
}

type SimpleUser struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type SimpleCommit struct {
	Author    SimpleUser `json:"author"`
	Committer SimpleUser `json:"committer"`
	Message   string     `json:"message"`
	URL       string     `json:"url"`
}

type Branch struct {
	Name   string     `json:"name"`
	Commit JSONCommit `json:"commit"`
}

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
