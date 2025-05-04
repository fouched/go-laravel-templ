package main

import (
	"github.com/fouched/rapidus"
	"log"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
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

	myMiddleware := &middleware.Middleware{
		App: rap,
	}

	myHandlers := &handlers.Handlers{App: rap}

	app := &application{
		App:        rap,
		Handlers:   myHandlers,
		Middleware: myMiddleware,
	}

	// set application routes to rapidus routes
	app.App.Routes = app.routes()

	// set the models
	app.Models = data.New(app.App.DB.Pool)
	myHandlers.Models = app.Models

	return app
}
