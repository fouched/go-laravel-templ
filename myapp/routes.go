package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *application) routes() *chi.Mux {
	// middleware comes before routes

	// routes
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
