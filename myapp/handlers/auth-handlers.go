package handlers

import (
	"github.com/fouched/rapidus/render"
	"myapp/views"
	"net/http"
)

func (h *Handlers) UserLoginGet(w http.ResponseWriter, r *http.Request) {

	t := views.Login()
	err := render.Template(w, r, t)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
