package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"git"
	"git/processor"
	"html/template"
	"jwts"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Page contains all data needed to Render the HTML files
type Page struct {
	Title,
	SinceDateString,
	GitClientID string
	RepoData []Repodata
	Settings git.Config
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

var jwtProvider jwts.Provider

// InitJWT intilialises the json web tocken Provider
func InitJWT() {
	jwtProvider = jwts.New([]byte(GlobalServerConfig.JWTSecret), jwts.Config{Method: jwt.SigningMethodHS256, TTL: 3600 * 24 * 3})
}

var page Page

var templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html", "login.html"))

func updatePageData(userCache *git.UserCache) {
	page.Title = TITLE
	page.GitClientID = "ea3fc9e6664643bd95b9"
	page.Authors = processor.GetCachedAuthors(userCache)
	page.Repos = processor.GetCachedRepos(userCache)
	page.Settings = userCache.Config
	page.SinceDateString = processor.GetCacheTimeString(userCache) //page.Settings.SinceTime.Format(time.RFC3339)[:10]
}

func gitHubSignIn(code string) (jwtBytes *jwt.Token, ID string, Error error) {

	authToken, err := git.GetAuthKeyFromGit(code, GlobalServerConfig.GitClientID, GlobalServerConfig.GitClientSecret)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	currentUser := git.GetUserFromToken(authToken)
	stringID := strconv.Itoa(currentUser.ID)
	if !user.Exists(stringID) {
		user.AddUser(stringID, authToken)
		processor.SaveCompleteConfig(user.GetUserCache(stringID), stringID)
	}

	token := jwtProvider.New()
	token.Claims["userID"] = stringID
	return token, stringID, nil
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userID string
	var err error
	var token *jwt.Token
	// if redirected from github login
	if code, ok := vars["githubLoginCode"]; ok {
		var newJwtBytes []byte
		token, userID, err = gitHubSignIn(code)
		if err != nil {
			// TODO
			log.Println(err)
		}
		newJwtBytes, err = jwtProvider.Sign(token)
		if err != nil {
			log.Println(err)
		}
		expiration := time.Now().Add(24 * 7 * time.Hour)
		cookie := http.Cookie{Name: "jwt", Value: string(newJwtBytes), Expires: expiration}
		http.SetCookie(w, &cookie)

	} else {
		userID, err = checkJWTandGetUserID(r)
		if err != nil {
			log.Println(err)
		}
	}

	if userID == "" {
		log.Println("Not logged in")
		updatePageData(&git.UserCache{})
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	if !user.ConfigLoaded(userID) {
		hasConfig := false
		for _, configFileName := range processor.GetSavedConfigs() {
			if configFileName == userID {
				hasConfig = true
				break
			}
		}
		if hasConfig {
			processor.LoadCompleteConfig(userID)
			user.SetConfigLoaded(userID, true)
		}

		// when config is loaded check JWT again against timestamp
		_, err = checkJWTandGetUserID(r)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}
	}

	// if cookie valid but AccessToken not found
	if _, ok := user.GetAccessToken(userID); !ok {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html", "login.html"))
	userCache := user.GetUserCache(userID)
	updatePageData(&userCache)

	templates.ExecuteTemplate(w, "commits.html", page)
}

// RefreshJWT refreshes the JWT cookie
func RefreshJWT(w http.ResponseWriter, r *http.Request) {
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	freshToken := jwtProvider.New()
	freshToken.Claims["userID"] = userID
	expiration := time.Now().Add(24 * 7 * time.Hour)
	freshTokenBytes, err := jwtProvider.Sign(freshToken)
	if err != nil {
		log.Fatalln(err)
	}
	updateIssuedAt(userID, freshToken.Claims["iatStr"].(string))
	cookie := http.Cookie{Name: "jwt", Value: string(freshTokenBytes), Expires: expiration}
	http.SetCookie(w, &cookie)
}

// CommitsShowJSON handler
func CommitsShowJSON(w http.ResponseWriter, r *http.Request) {
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		return
	}
	if !user.Exists(userID) {
		return
	}
	cacheTimeBefore := user.GetUserCache(userID).CacheTime
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Page-has-updates", "false")
	vars := map[string]string{"commit": r.FormValue("commit"), "author": r.FormValue("author"), "repo": r.FormValue("repo"), "date": r.FormValue("date")}
	query := getQueryFromVars(vars)

	queryResult := git.Commits{}
	commitChanel := make(chan git.Commit)
	go processor.SendCommits(userID, query, commitChanel)
	for commit := range commitChanel {
		queryResult = append(queryResult, commit)
	}
	sort.Sort(queryResult)

	cacheTimeAfterQuery := user.GetUserCache(userID).CacheTime
	if cacheTimeAfterQuery.Before(cacheTimeBefore) {
		w.Header().Set("Page-has-updates", "true")
	}

	//For paging. On empty or error: no paging
	page, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("perPage"))
	var maxIndex, minIndex int

	if page != 0 {
		// if not specified 30 commits per page
		if perPage == 0 {
			perPage = 30
		}
		maxIndex = page*perPage - 1
		minIndex = page*perPage - perPage
	}

	currentIndex := -1
	commitData := []Buttondata{}
	for _, com := range queryResult {
		currentIndex++
		formatedDate := com.Time.Format(time.RFC822)[:10]
		if page != 0 {
			if currentIndex < minIndex {
				continue
			}
			if currentIndex > maxIndex {
				break
			}
		}
		commitData = append(commitData, Buttondata{com.Comment, com.Sha, formatedDate, com.Repo + "/" + com.Branch, com.Time.UnixNano()})
	}
	if err := json.NewEncoder(w).Encode(commitData); err != nil {
		panic(err)
	}
}

// AuthorsShowJSON handler
func AuthorsShowJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		return
	}
	userCache := user.GetUserCache(userID)
	if err := json.NewEncoder(w).Encode(processor.GetCachedAuthors(&userCache)); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

// AuthorsShow handler
func AuthorsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
	}
	if !user.Exists(userID) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}
	userCache := user.GetUserCache(userID)
	authorButtons := []Buttondata{}
	for _, author := range processor.GetCachedAuthors(&userCache) {
		authorButtons = append(authorButtons, Buttondata{author, author, "", "", 0})
	}
	updatePageData(&userCache)
	templates.ExecuteTemplate(w, "authors.html", page)
}

// LoginHTML handler
func LoginHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	updatePageData(&git.UserCache{})
	templates = template.Must(template.ParseFiles("headAndNavbar.html", "scripts.html", "login.html"))
	templates.ExecuteTemplate(w, "login.html", page)
}

// SettingsShow handler
func SettingsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
	}
	if !user.Exists(userID) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	userCache := user.GetUserCache(userID)
	updatePageData(&userCache)
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
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		return
	}
	processor.SetConfig(userID, config)
	userCache := user.GetUserCache(userID)
	processor.SaveCompleteConfig(userCache, userID)
	updatePageData(&userCache)
	templates.ExecuteTemplate(w, "settings.html", page)
}

// ReposShowHTML handler
func ReposShowHTML(w http.ResponseWriter, r *http.Request) {
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
	}
	if !user.Exists(userID) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}
	userCache := user.GetUserCache(userID)
	w.Header().Set("Content-type", "text/html")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html", "login.html"))
	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	repos, err := processor.GetCachedRepoObjects(userID)
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
	updatePageData(&userCache)
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

	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		return
	}

	processor.SetRepoBranch(userID, repo, branch)
}

// ReposShow handler
func ReposShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html", "repositories.html", "settings.html", "authors.html", "scripts.html"))
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		return
	}
	userCache := user.GetUserCache(userID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedRepos(&userCache)); err != nil {
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
	userID, err := checkJWTandGetUserID(r)
	if err != nil {
		log.Println(err)
		return
	}
	userCache := user.GetUserCache(userID)
	if err := json.NewEncoder(w).Encode(processor.GetSingleCommit(userCache, sha)); err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	query.Commit = vars["commit"]
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

func checkJWTandGetUserID(r *http.Request) (string, error) {
	var userID string
	JWT, err := jwtProvider.Get(r)
	if err != nil {
		// TODO no one logged in
		return "", err
	}
	userID = JWT.Claims["userID"].(string)
	if userID == "" {
		return "", errors.New("Error: Empty userID in JWT!")
	}
	if JWT.Claims["iatStr"] == nil {
		return "", errors.New("Error: JWT valid but no claim 'iatStr'!")
	}
	issuedAt := JWT.Claims["iatStr"].(string)
	// if user isn't created yet there should be no previous issuedAt time
	if user.Exists(userID) {
		userConfig := user.GetUserCache(userID).Config
		if issuedAt < userConfig.JWTissuedAt {
			return "", errors.New("Error: JWT valid but expired! (not the latest one)")
		}
	}

	return userID, nil
}

func updateIssuedAt(userID string, issuedAt string) {
	cache := user.GetUserCache(userID)
	config := cache.Config
	config.JWTissuedAt = issuedAt
	cache.Config = config
	user.SetUserCache(userID, cache)
	processor.SaveCompleteConfig(cache, userID)
}
