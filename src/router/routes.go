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
		nil,
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
		"/repos",
		nil,
		ReposShow,
	},
	Route{
		"CommitShow",
		"GET",
		"/commits",
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
		"SetConfig",
		"GET",
		"/setConfig={jsonString:(?s).*}",
		nil,
		SetConfig,
	},
	Route{
		"Config",
		"GET",
		"/config",
		nil,
		GetConfig,
	},
}
