package router

import (
	"encoding/json"
	"fmt"
	"git/processor"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
	fmt.Fprintln(w, "")
}

func shutdownCHM(w http.ResponseWriter, r *http.Request) {
	defer os.Exit(0)
}

func Log(w http.ResponseWriter, r *http.Request) {
	w.Write(LogBuffer.Bytes())
}

func AuthorsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedAuthors()); err != nil {
		panic(err)
	}
}
func ReposShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedRepos()); err != nil {
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
