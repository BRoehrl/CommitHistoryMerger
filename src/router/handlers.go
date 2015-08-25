package router

import (
	"encoding/json"
	"fmt"
	"git"
	"git/processor"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title      string
	Buttondata []Buttondata
	RepoData   []Repodata
	Settings   git.Config
}
type Repodata struct {
	Name       string
	Branches   []string
	NrBranches int
}
type Buttondata struct {
	Name,
	Id,
	DateString,
	Repository string
}

const (
	TITLE = "CHM"
)

var templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	//	threeMonthAgo := time.Now().AddDate(0, -3, 0)
	//	query := processor.Query{Since: threeMonthAgo}
	vars := mux.Vars(r)

	query := getQueryFromVars(vars)

	queryResult := processor.GetCommits(query)

	commitData := []Buttondata{}
	for _, com := range queryResult {
		formatedDate := com.Time.Format(time.RFC822)[:10]
		commitData = append(commitData, Buttondata{com.Comment, com.Sha, formatedDate, com.Repo + "/" + com.Branch})
	}
	templates.ExecuteTemplate(w, "commits.html", Page{Title: TITLE, Buttondata: commitData}) //Page{Title: "Home"})
}

func shutdownCHM(w http.ResponseWriter, r *http.Request) {
	defer os.Exit(0)
}

func Log(w http.ResponseWriter, r *http.Request) {
	w.Write(LogBuffer.Bytes())
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
	templates.ExecuteTemplate(w, "commits.html", Page{Title: TITLE, Buttondata: authorButtons})
}

func SettingsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	templates.ExecuteTemplate(w, "settings.html", Page{Title: TITLE, Settings: git.GetConfig()})
}

func SettingsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}

	config := getConfigFromForm(r.Form)
	fmt.Println(config)
	processor.SetConfig(config)
	templates.ExecuteTemplate(w, "settings.html", Page{Title: TITLE, Settings: git.GetConfig()})
}

func ReposShowHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
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
	templates.ExecuteTemplate(w, "repositories.html", Page{Title: TITLE, RepoData: repodata})
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

	} else {
		if d, err := time.Parse("2006-01-02", git.GetConfig().SinceTime); err == nil {
			query.Since = d
		}
	}
	return query
}

func getConfigFromForm(form url.Values) git.Config {
	config := git.Config{}
	config.GitUrl = form.Get("baseUrl")
	config.BaseOrganisation = form.Get("baseOrg")
	config.GitAuthkey = form.Get("authKey")
	config.SinceTime = form.Get("sinceTime")
	config.MaxRepos, _ = strconv.Atoi(form.Get("maxRepos"))
	config.MaxBranches, _ = strconv.Atoi(form.Get("maxBranches"))

	return config
}
