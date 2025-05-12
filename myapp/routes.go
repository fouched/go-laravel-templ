package main

import (
	"fmt"
	"github.com/fouched/rapidus/mailer"
	"github.com/go-chi/chi/v5"
	"myapp/data"
	"net/http"
	"strconv"
)

func (a *application) routes() *chi.Mux {
	// middleware comes before routes
	a.use(a.Middleware.CheckRemember)

	// routes
	//a.App.Routes.Get("/", a.Handlers.Home)
	//using the convenience wrapper below
	a.get("/", a.Handlers.Home)
	a.get("/sessions", a.Handlers.SessionTest)

	a.get("/users/login", a.Handlers.UserLoginGet)
	a.post("/users/login", a.Handlers.UserLoginPost)
	a.get("/users/logout", a.Handlers.LogOut)
	a.get("/users/forgot-password", a.Handlers.ForgotGet)
	a.post("/users/forgot-password", a.Handlers.ForgotPost)
	a.get("/users/reset-password", a.Handlers.ResetPasswordForm)

	a.get("/form", a.Handlers.Form)
	a.post("/form", a.Handlers.PostForm)

	a.get("/json", a.Handlers.JSON)
	a.get("/xml", a.Handlers.XML)
	a.get("/download-file", a.Handlers.DownloadFile)

	a.get("/crypto", a.Handlers.TestCrypto)

	a.get("/cache/demo", a.Handlers.CacheDemo)
	a.post("/cache/save", a.Handlers.CacheSave)
	a.post("/cache/get", a.Handlers.CacheGet)
	a.post("/cache/delete", a.Handlers.CacheDelete)
	a.post("/cache/empty", a.Handlers.CacheEmpty)

	a.get("/cache/redis/demo", a.Handlers.CacheDemoRedis)
	a.post("/cache/redis/save", a.Handlers.CacheSaveRedis)
	a.post("/cache/redis/get", a.Handlers.CacheGetRedis)
	a.post("/cache/redis/delete", a.Handlers.CacheDeleteRedis)
	a.post("/cache/redis/empty", a.Handlers.CacheEmptyRedis)

	a.get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "me@here.com",
			To:          "you@there.com",
			Subject:     "Test using channel",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		//a.App.Mail.Jobs <- msg
		//res := <-a.App.Mail.Results
		//if res.Error != nil {
		//	a.App.ErrorLog.Println(res.Error)
		//}

		msg.Subject = "Test using direct call"
		err := a.App.Mail.SendSMTPMessage(msg)
		if err != nil {
			a.App.ErrorLog.Println(err)
		}
	})

	a.get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "Fouche",
			LastName:  "du Preez",
			Email:     "me@here.com",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d: %s", id, u.FirstName)
	})

	a.get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		for _, x := range users {
			fmt.Fprintf(w, x.LastName)
		}
	})

	a.get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})

	a.get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		u.LastName = a.App.RandomString(10)

		// test validator
		validator := a.App.Validator(nil)
		u.LastName = ""
		u.Validate(validator)

		if !validator.Valid() {
			fmt.Fprintf(w, "failed validation")
			return
		}

		err = u.Update(*u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
