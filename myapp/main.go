package main

import (
	"github.com/fouched/rapidus"
	"myapp/data"
	"myapp/handlers"
)

type application struct {
	App      *rapidus.Rapidus
	Handlers *handlers.Handlers
	Models   data.Models
}

func main() {
	r := initApplication()
	r.App.ListenAndServe()
}
