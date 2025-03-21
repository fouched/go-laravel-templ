package main

import (
	"github.com/fouched/rapidus"
	"myapp/handlers"
)

type application struct {
	App      *rapidus.Rapidus
	Handlers *handlers.Handlers
}

func main() {
	r := initApplication()
	r.App.ListenAndServe()
}
