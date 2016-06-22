package router

import (
	"net/http"
)

// Route contains all information to create a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Queries     [][]string
	HandlerFunc http.HandlerFunc
}

// Routes is a slice of Route structs
type Routes []Route

var routes = Routes{
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
			[]string{"code", "{githubLoginCode}"},
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
		"RefreshJWT",
		"GET",
		"/refresh_token",
		nil,
		RefreshJWT,
	},
	Route{
		"SettingsShow",
		"GET",
		"/settings",
		nil,
		SettingsShow,
	},
	Route{
		"SettingsPost",
		"POST",
		"/settings",
		nil,
		SettingsPost,
	},
	Route{
		"ReposShow",
		"GET",
		"/json/repos",
		nil,
		ReposShow,
	},
	Route{
		"ReposShowHTML",
		"GET",
		"/repositories",
		nil,
		ReposShowHTML,
	},
	Route{
		"Login",
		"GET",
		"/login",
		nil,
		LoginHTML,
	},
	Route{
		"RepoBranchChange",
		"POST",
		"/repositories",
		nil,
		RepoBranchChange,
	},
	Route{
		"CommitsShowJSON",
		"POST",
		"/commits",
		nil,
		CommitsShowJSON,
	},
	Route{
		"SingleCommit",
		"GET",
		"/json/commits/{sha}",
		nil,
		ShowSingleCommit,
	},
}
