package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		if route.Queries == nil {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		} else {
			for _, query := range route.Queries {
				router.
					Queries(query...).
					Methods(route.Method).
					Path(route.Pattern).
					Name(route.Name).
					Handler(handler)
			}
		}

	}
	router.HandleFunc("/socket", SocketHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	return router
}
