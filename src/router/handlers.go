package router

import (
	"encoding/json"
	"fmt"
	"git/processor"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title       string
	Buttondata []Buttondata
}
type Buttondata struct {
	Name string
	Id string
}

const (
	TITLE = "CHM"
)

var templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	fmt.Println("Formvalue:",r.FormValue("user"))
	fmt.Println("Formvalue:",r.FormValue("pw"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	threeMonthAgo := time.Now().AddDate(0, -3, 0)
	query := processor.Query{Since: threeMonthAgo}
	queryResult := processor.GetCommits(query)
	
	commitData := []Buttondata{}
	for _, com := range(queryResult){
		formatedName := com.Time.Format(time.RFC822)[:7] + com.Comment
		commitData = append(commitData, Buttondata{formatedName, com.Sha})
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
	for _, author := range(processor.GetCachedAuthors()){
		authorButtons = append(authorButtons, Buttondata{author, author})
	}
	templates.ExecuteTemplate(w, "commits.html", Page{Title: TITLE, Buttondata: authorButtons})
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

	threeMonthAgo := time.Now().AddDate(0, -3, 0)
	query := processor.Query{Since: threeMonthAgo}

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

	queryResult := processor.GetCommits(query)

	if err := json.NewEncoder(w).Encode(queryResult); err != nil {
		panic(err)
	}
}

func SetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString := vars["jsonString"]
	jsonString = strings.Replace(jsonString, "+", "", -1)
	fmt.Fprintln(w, jsonString)
	if err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&[]string{}); err != nil {
		fmt.Fprintln(w, err)
		return
	}
}
func GetConfig(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode("Test"); err != nil {
		panic(err)
	}
}

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}
