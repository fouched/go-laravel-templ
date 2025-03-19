package main

import (
	"github.com/fouched/rapidus"
	"log"
	"os"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init celeritas
	rap := &rapidus.Rapidus{}
	err = rap.New(path)
	if err != nil {
		log.Fatal(err)
	}

	rap.AppName = "myapp"
	rap.Debug = true

	app := &application{
		App: rap,
	}

	return app
}
