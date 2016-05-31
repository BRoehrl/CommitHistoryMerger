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
		"SettingsShow",
		"GET",
		"/settings",
		nil,
		SettingsShow,
	},
	Route{
		"SettingsShow",
		"POST",
		"/settings",
		nil,
		SettingsPost,
	},
	Route{
		"ProfileSave",
		"POST",
		"/config/save/{name}",
		nil,
		SaveProfile,
	},
	Route{
		"ProfileLoad",
		"GET",
		"/config/load/{name}",
		nil,
		LoadProfile,
	},
	Route{
		"ReposShow",
		"GET",
		"/json/repos",
		nil,
		ReposShow,
	},
	Route{
		"RepositoryHTML",
		"GET",
		"/repositories",
		nil,
		ReposShowHTML,
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
