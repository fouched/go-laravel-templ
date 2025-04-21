package handlers

import (
	"github.com/fouched/rapidus"
	"myapp/data"
	"myapp/views"
	"net/http"
)

type Handlers struct {
	App    *rapidus.Rapidus
	Models data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionGetInt(r.Context(), "userID")
	isAuthenticated := userID != 0

	h.render(w, r, views.Home(isAuthenticated))
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	myData := "bar"
	h.sessionPut(r.Context(), "foo", myData)
	myValue := h.sessionGetString(r.Context(), "foo")

	h.render(w, r, views.Sessions(myValue))
}
