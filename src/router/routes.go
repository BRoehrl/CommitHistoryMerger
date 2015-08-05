package router

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Queries     [][]string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"shutdownCHM",
		"GET",
		"/shutdown",
		nil,
		shutdownCHM,
	},
	Route{
		"log",
		"GET",
		"/log",
		nil,
		Log,
	},
	Route{
		"Index",
		"GET",
		"/",
		[][]string{
			[]string{"author", "{author}", "repo", "{repo}", "since", "{date}"},
			[]string{"author", "{author}", "repo", "{repo}"},
			[]string{"author", "{author}", "since", "{date}"},
			[]string{"repo", "{repo}", "since", "{date}"},
			[]string{"author", "{author}"},
			[]string{"repo", "{repo}"},
			[]string{"since", "{date}"},
			nil,
		},
		Index,
	},
	Route{
		"AuthorsShow",
		"GET",
		"/authors",
		nil,
		AuthorsShow,
	},
	Route{
		"ReposShow",
		"GET",
		"/json/repos",
		nil,
		ReposShow,
	},
	Route{
		"CommitShow",
		"GET",
		"/json/commits",
		[][]string{
			[]string{"author", "{author}", "repo", "{repo}", "since", "{date}"},
			[]string{"author", "{author}", "repo", "{repo}"},
			[]string{"author", "{author}", "since", "{date}"},
			[]string{"repo", "{repo}", "since", "{date}"},
			[]string{"author", "{author}"},
			[]string{"repo", "{repo}"},
			[]string{"since", "{date}"},
			nil,
		},
		CommitShow,
	},
	Route{
		"SingleCommit",
		"GET",
		"/json/commits/{sha}",
		nil,
		ShowSingleCommit,
	},
	Route{
		"SetConfig",
		"GET",
		"/json/setConfig={jsonString:(?s).*}",
		nil,
		SetConfig,
	},
	Route{
		"Config",
		"GET",
		"/json/config",
		nil,
		GetConfig,
	},
}
