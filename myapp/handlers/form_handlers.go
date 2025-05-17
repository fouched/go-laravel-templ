package handlers

import (
	"fmt"
	"github.com/fouched/rapidus"
	"myapp/data"
	"myapp/views"
	"net/http"
)

func (h *Handlers) Form(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, views.Form(data.User{}, h.App.Validator(nil)))
}

func (h *Handlers) PostForm(w http.ResponseWriter, r *http.Request) {
	// all form posts must be parsed
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	validator := h.App.Validator(nil)

	validator.Required(
		rapidus.Field{
			Name:  "first_name",
			Label: "First Name",
			Value: r.Form.Get("first_name"),
		},
		rapidus.Field{
			Name:  "last_name",
			Label: "Last Name",
			Value: r.Form.Get("last_name"),
		},
		rapidus.Field{
			Name:  "email",
			Label: "Email",
			Value: r.Form.Get("email"),
		},
	)

	validator.IsEmail(rapidus.Field{
		Name:  "email",
		Label: "Email",
		Value: r.Form.Get("email"),
	})

	validator.Check(len(r.Form.Get("first_name")) > 1, "first_name", "First Name must be at least two characters")
	validator.Check(len(r.Form.Get("last_name")) > 1, "last_name", "Last Name must be at least two characters")

	if validator.Valid() {
		fmt.Fprint(w, "valid data")
	} else {
		var user data.User
		user.FirstName = r.Form.Get("first_name")
		user.LastName = r.Form.Get("last_name")
		user.Email = r.Form.Get("email")

		h.render(w, r, views.Form(user, validator))
	}
}
