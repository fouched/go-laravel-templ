package main

import (
	"github.com/fouched/rapidus"
	"log"
	"myapp/data"
	"myapp/handlers"
	"os"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init rapidus
	rap := &rapidus.Rapidus{}
	err = rap.New(path)
	if err != nil {
		log.Fatal(err)
	}

	rap.AppName = "myapp"

	myHandlers := &handlers.Handlers{App: rap}

	app := &application{
		App:      rap,
		Handlers: myHandlers,
	}

	// set application routes to rapidus routes
	app.App.Routes = app.routes()
	app.Models = data.New(app.App.DB.Pool)

	return app
}
