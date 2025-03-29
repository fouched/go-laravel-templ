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

	h.App.Session.Put(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) LogOut(w http.ResponseWriter, r *http.Request) {

	h.App.Session.RenewToken(r.Context())
	h.App.Session.Remove(r.Context(), "userID")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
