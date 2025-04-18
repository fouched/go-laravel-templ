package middleware

import (
	"github.com/fouched/rapidus"
	"myapp/data"
)

type Middleware struct {
	App    *rapidus.Rapidus
	Models data.Models
}
