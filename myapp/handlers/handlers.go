package handlers

import (
	"github.com/fouched/rapidus"
	"myapp/data"
	"myapp/views"
	"net/http"
	"time"
)

type Handlers struct {
	App    *rapidus.Rapidus
	Models data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	defer h.App.LoadTime(time.Now())

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

func (h *Handlers) JSON(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID      int      `json:"ID"`
		Name    string   `json:"name"`
		Hobbies []string `json:"hobbies"`
	}

	payload.ID = 10
	payload.Name = "Jack Jones"
	payload.Hobbies = []string{"karate", "tennis", "programming"}

	err := h.App.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) XML(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		ID      int      `xml:"ID"`
		Name    string   `xml:"name"`
		Hobbies []string `xml:"hobbies>hobby"`
	}

	var payload Payload
	payload.ID = 10
	payload.Name = "John Smith"
	payload.Hobbies = []string{"karate", "tennis", "programming"}

	if err := h.App.WriteXML(w, http.StatusOK, payload); err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	h.App.DownloadStaticFile(w, r, "./public/images/rapidus.jpg", "rapidus.jpg")
}

func (h *Handlers) TestCrypto(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello, world"

	encrypted, err := h.encrypt(plainText)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Error500(w)
	}

	decrypted, err := h.decrypt(encrypted)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Error500(w)
	}

	h.render(w, r, views.Crypto(plainText, encrypted, decrypted))
}
