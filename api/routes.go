/*
routes seperates routes from the main application logic, define routes or paths in
defineRoutes().
*/

package api

import (
	"net/http"

	handlers "github.com/deejcoder/spidernet-api/handlers"
	helpers "github.com/deejcoder/spidernet-api/helpers"
	"github.com/gorilla/mux"
)

// Route defines a Route
type Route struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
	Methods string
}

// BuildRouter creates a new Router and adds all defined routes to it
func BuildRouter() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	routes := defineRoutes()
	for _, route := range routes {
		api.HandleFunc(route.Path, route.Handler).Methods(route.Methods)
	}
	return router
}

func defineRoutes() []Route {
	return []Route{
		{Path: "/", Handler: helpers.RequireAuth(handlers.Index), Methods: "GET"},
		{Path: "/auth/login", Handler: handlers.LoginUser, Methods: "POST"},
		{Path: "/auth/validate", Handler: handlers.ValidateUser, Methods: "GET"},
		{Path: "/servers/", Handler: handlers.GetServers, Methods: "GET"},
		{Path: "/servers/search", Handler: handlers.SearchServers, Methods: "GET"},
		{Path: "/servers/create", Handler: handlers.CreateServer, Methods: "PUT"},
		{Path: "/servers/delete", Handler: handlers.DeleteServer, Methods: "DELETE"},
		{Path: "/servers/update", Handler: handlers.UpdateServer, Methods: "POST"},
	}
}
