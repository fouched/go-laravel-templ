package handlers

import (
	"myapp/views"
	"net/http"
)

func (h *Handlers) UserLoginGet(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, views.Login())
}

func (h *Handlers) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err.Error())
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		h.App.ErrorLog.Println(err.Error())
		return
	}

	matches, err := user.PasswordMatches(password)
	if err != nil {
		h.App.ErrorLog.Println(err.Error())
		return
	}
	if !matches {
		h.App.ErrorLog.Println("Invalid password")
		return
	}

	h.sessionPut(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) LogOut(w http.ResponseWriter, r *http.Request) {
	_ = h.sessionRenew(r.Context())
	h.sessionRemove(r.Context(), "userID")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
