package git

import (
	"time"
)

type Repo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	Language string `json:"language"`
}

type Repos []Repo

type GitUser struct {
	Id      int    `json:"id"`
	Name    string `json:"login"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

type JsonCommit struct {
	Sha          string       `json:"sha"`
	ActualCommit SimpleCommit `json:"commit"`
	Url          string       `json:"url"`
	HtmlUrl      string       `json:"html_url"`
	Author       GitUser      `json:"author"`
	Committer    GitUser      `json:"committer"`
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
	Url       string     `json:"url"`
}

type SortableCommits []JsonCommit

//type SortableCommits []Commit

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
