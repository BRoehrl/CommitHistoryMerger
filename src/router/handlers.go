package router

import (
	"encoding/json"
	"fmt"
	"git"
	"git/processor"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Page contains all data needed to Render the HTML files
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

// Repodata contains data for rendering the repository list
type Repodata struct {
	Name       string
	Branches   []string
	NrBranches int
}

// Buttondata contains data for rendering the commit button list
type Buttondata struct {
	Name,
	ID,
	DateString,
	Repository string
	NanoTime int64
}

const (
	// TITLE is the title of the index page
	TITLE = "CHF"
)

var page Page

var templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))

func updatePageData() {
	page.Title = TITLE
	page.Profiles = processor.GetSavedConfigs()
	page.Authors = processor.GetCachedAuthors()
	page.Repos = processor.GetCachedRepos()
	page.Settings = git.GetConfig()
	page.SinceDateString = processor.GetCacheTimeString() //page.Settings.SinceTime.Format(time.RFC3339)[:10]
	page.ActiveProfile = processor.LoadedConfig
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))

	vars := mux.Vars(r)

	query := getQueryFromVars(vars)

	//queryResult := processor.GetCommits(query)
	//Experimental: Concurrent implementation with chanels
	queryResult := processor.Commits{}
	commitChanel := make(chan git.Commit)
	go processor.SendCommits(query, commitChanel)
	for commit := range commitChanel {
		queryResult = append(queryResult, commit)
	}
	sort.Sort(queryResult)
	//log.Println(len(queryResult))

	commitData := []Buttondata{}
	for _, com := range queryResult {
		formatedDate := com.Time.Format(time.RFC822)[:10]
		commitData = append(commitData, Buttondata{com.Comment, com.Sha, formatedDate, com.Repo + "/" + com.Branch, com.Time.UnixNano()})
	}
	page.Buttondata = commitData
	updatePageData()
	templates.ExecuteTemplate(w, "commits.html", page)
}

// AuthorsShowJSON handler
func AuthorsShowJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedAuthors()); err != nil {
		panic(err)
	}
}

// AuthorsShow handler
func AuthorsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	authorButtons := []Buttondata{}
	for _, author := range processor.GetCachedAuthors() {
		authorButtons = append(authorButtons, Buttondata{author, author, "", "", 0})
	}
	updatePageData()
	templates.ExecuteTemplate(w, "authors.html", page)
}

// SettingsShow handler
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

// SettingsPost handler
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

// SaveProfile handler
func SaveProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	vars := mux.Vars(r)
	processor.SaveCompleteConfig(vars["name"])
}

// LoadProfile handler
func LoadProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	vars := mux.Vars(r)
	processor.LoadCompleteConfig(vars["name"])
}

// ReposShowHTML handler
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

// RepoBranchChange handler
func RepoBranchChange(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	repo := r.FormValue("repo")
	branch := r.FormValue("branch")
	processor.SetRepoBranch(repo, branch)
}

// ReposShow handler
func ReposShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedRepos()); err != nil {
		panic(err)
	}
}

// ShowSingleCommit handler
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// SocketHandler handler
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		jVars := git.JSONVars{}
		err := conn.ReadJSON(&jVars)
		if err != nil {
			//log.Println("1", err)
			return
		}
		vars := map[string]string{"author": jVars.Author, "repo": jVars.Repo, "date": jVars.Querydate}
		query := getQueryFromVars(vars)
		commitChanel := make(chan git.Commit)
		go processor.SendCommits(query, commitChanel)
		buttonBuffer := []Buttondata{}
		for com := range commitChanel {
			formatedDate := com.Time.Format(time.RFC822)[:10]
			bdata := Buttondata{com.Comment, com.Sha, formatedDate, com.Repo + "/" + com.Branch, com.Time.UnixNano()}
			buttonBuffer = append(buttonBuffer, bdata)
			if len(buttonBuffer) > 9 {
				err = conn.WriteJSON(buttonBuffer)
				if err != nil {
					log.Println("2", err)
					return
				}
				buttonBuffer = []Buttondata{}
			}
		}
		if err != nil {
			log.Println("3", err)
			return
		}
	}
}

// CommitShow handler
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

	authors := vars["author"]
	if authors != "" {
		if form, err := url.QueryUnescape(authors); err == nil {
			authors = form
		}
		query.Authors = strings.Split(authors, ";")
	}

	repos := vars["repo"]
	if repos != "" {
		if form, err := url.QueryUnescape(repos); err == nil {
			repos = form
		}
		query.Repos = strings.Split(repos, ";")
	}

	since := vars["date"]
	if since != "" {
		if d, err := time.Parse(time.RFC3339, since); err == nil {
			query.Since = d
		}
	}

	query.UseRegex = true

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
