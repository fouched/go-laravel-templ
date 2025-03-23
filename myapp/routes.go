package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *application) routes() *chi.Mux {
	// middleware comes before routes

	// routes
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)

	// test func for db connectivity
	a.App.Routes.Get("/test-database", func(w http.ResponseWriter, r *http.Request) {
		query := "select id, first_name from users where id = 1"

		row := a.App.DB.Pool.QueryRowContext(r.Context(), query)

		var id int
		var name string
		err := row.Scan(&id, &name)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		_, _ = fmt.Fprintf(w, "%d %s", id, name)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
