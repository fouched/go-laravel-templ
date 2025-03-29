package handlers

import (
	"github.com/fouched/rapidus"
	"github.com/fouched/rapidus/render"
	"myapp/data"
	"myapp/views"
	"net/http"
)

type Handlers struct {
	App    *rapidus.Rapidus
	Models data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {

	userID := h.App.Session.GetInt(r.Context(), "userID")
	isAuthenticated := userID != 0
	t := views.Home(isAuthenticated)
	err := render.Template(w, r, t)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	myData := "bar"
	h.App.Session.Put(r.Context(), "foo", myData)

	myValue := h.App.Session.GetString(r.Context(), "foo")

	t := views.Sessions(myValue)
	err := render.Template(w, r, t)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
