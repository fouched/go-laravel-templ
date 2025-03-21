package handlers

import (
	"github.com/fouched/rapidus"
	"github.com/fouched/rapidus/render"
	"myapp/views"
	"net/http"
)

type Handlers struct {
	App *rapidus.Rapidus
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	t := views.Home()
	err := render.Template(w, r, t)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
