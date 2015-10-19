package router

import (
	"encoding/json"
	"fmt"
	"git"
	"git/processor"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title,
	SinceDateString,
	ActiveProfile string
	Buttondata []Buttondata
	RepoData   []Repodata
	Settings   git.Config
	Profiles,
	Authors,
	Repos []string
}
type Repodata struct {
	Name       string
	Branches   []string
	NrBranches int
}
type Buttondata struct {
	Name,
	ID,
	DateString,
	Repository string
}

const (
	TITLE = "CHM"
)

var page Page

var templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))

func updatePageData() {
	page.Title = TITLE
	page.Profiles = processor.GetSavedConfigs()
	page.Authors = processor.GetCachedAuthors()
	page.Repos = processor.GetCachedRepos()
	page.Settings = git.GetConfig()
	page.SinceDateString = page.Settings.SinceTime.Format(time.RFC3339)[:10]
	page.ActiveProfile = processor.LoadedConfig
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))

	vars := mux.Vars(r)

	query := getQueryFromVars(vars)

	queryResult := processor.GetCommits(query)

	commitData := []Buttondata{}
	for _, com := range queryResult {
		formatedDate := com.Time.Format(time.RFC822)[:10]
		commitData = append(commitData, Buttondata{com.Comment, com.Sha, formatedDate, com.Repo + "/" + com.Branch})
	}
	page.Buttondata = commitData
	updatePageData()
	templates.ExecuteTemplate(w, "commits.html", page)
}

func AuthorsShowJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedAuthors()); err != nil {
		panic(err)
	}
}

func AuthorsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	authorButtons := []Buttondata{}
	for _, author := range processor.GetCachedAuthors() {
		authorButtons = append(authorButtons, Buttondata{author, author, "", ""})
	}
	updatePageData()
	templates.ExecuteTemplate(w, "authors.html", page)
}

func SettingsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	updatePageData()
	templates.ExecuteTemplate(w, "settings.html", page)
}

func SettingsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}

	config := getConfigFromForm(r.Form)
	processor.SetConfig(config)
	updatePageData()
	templates.ExecuteTemplate(w, "settings.html", page)
}

func SaveProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	vars := mux.Vars(r)
	processor.SaveCompleteConfig(vars["name"])
}

func LoadProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	vars := mux.Vars(r)
	processor.LoadCompleteConfig(vars["name"])
}

func ReposShowHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	repos, err := processor.GetCachedRepoObjects()
	if err != nil {
		panic(err)
	}
	repodata := []Repodata{}
	for _, repo := range repos {
		branches := []string{repo.SelectedBranch}
		for branch := range repo.Branches {
			if branch != repo.SelectedBranch {
				branches = append(branches, branch)
			}
		}
		repodata = append(repodata, Repodata{repo.Name, branches, len(branches)})

	}
	updatePageData()
	page.RepoData = repodata
	templates.ExecuteTemplate(w, "repositories.html", page)
}
func RepoBranchChange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repo, ok := vars["repo"]
	if !ok { // TODO send error
		return
	}
	branch, ok := vars["branch"]
	if !ok { // TODO send error
		return
	}
	processor.SetRepoBranch(repo, branch)
}

func ReposShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedRepos()); err != nil {
		panic(err)
	}
}

func ShowSingleCommit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)
	sha, ok := vars["sha"]
	if !ok {
		// TODO
	}
	if err := json.NewEncoder(w).Encode(processor.GetSingleCommit(sha)); err != nil {
		panic(err)
	}
}

func CommitShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)

	query := getQueryFromVars(vars)

	queryResult := processor.GetCommits(query)

	if err := json.NewEncoder(w).Encode(queryResult); err != nil {
		panic(err)
	}
}

func getQueryFromVars(vars map[string]string) processor.Query {

	query := processor.Query{}

	authors, ok := vars["author"]
	if ok {
		if form, err := url.QueryUnescape(authors); err == nil {
			authors = form
		}
		query.Authors = strings.Split(authors, ";")
	}

	repos, ok := vars["repo"]
	if ok {
		if form, err := url.QueryUnescape(repos); err == nil {
			repos = form
		}
		query.Repos = strings.Split(repos, ";")
	}

	since, ok := vars["date"]
	if ok {
		if d, err := time.Parse(time.RFC3339, since); err == nil {
			query.Since = d
		}
	}
	return query
}

func getConfigFromForm(form url.Values) git.Config {
	config := git.Config{}
	config.GitURL = form.Get("baseURL")
	config.BaseOrganisation = form.Get("baseOrg")
	config.GitAuthkey = form.Get("authKey")
	config.MiscDefaultBranch = form.Get("defaultBranch")
	if d, err := time.Parse("2006-01-02", form.Get("sinceTime")); err == nil {
		// add a day to include commits on day 'since'
		config.SinceTime = d.AddDate(0, 0, 1)
	} else {
		config.SinceTime = time.Time{}
	}
	config.MaxRepos, _ = strconv.Atoi(form.Get("maxRepos"))
	config.MaxBranches, _ = strconv.Atoi(form.Get("maxBranches"))

	return config
}
