package handlers

import (
	"myapp/views"
	"net/http"
)

func (h *Handlers) CacheDemo(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, views.CachePage())
}

func (h *Handlers) CacheSave(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
	}

	type userInput struct {
		Key   string
		Value string
		CSRF  string
	}

	input := userInput{
		Key:   r.Form.Get("cache_name"),
		Value: r.Form.Get("cache_value"),
	}

	msg := "Saved in cache"

	err = h.App.RedisClient.Set(r.Context(), input.Key, input.Value, 0).Err()
	if err != nil {
		h.App.ErrorLog.Println(err)
		msg = "Could not save to cache"
	}

	h.render(w, r, views.CacheSave(msg, err))
}

func (h *Handlers) CacheGet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
	}

	key := r.Form.Get("cache_get")
	msg := "Could not get entry from cache"

	val, err := h.App.RedisClient.Get(r.Context(), key).Result()
	if err == nil {
		msg = val
	} else {
		h.App.ErrorLog.Println(err)
	}

	h.render(w, r, views.CacheGet(msg, err))
}

func (h *Handlers) CacheDelete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
	}

	key := r.Form.Get("cache_delete")
	msg := "Deleted cache key: " + key

	err = h.App.RedisClient.Del(r.Context(), key).Err()
	if err != nil {
		h.App.ErrorLog.Println(err)
		msg = "Could not delete key: " + key
	}

	h.render(w, r, views.CacheDelete(msg, err))
}

func (h *Handlers) CacheEmpty(w http.ResponseWriter, r *http.Request) {
	msg := "Emptied cache"
	err := h.App.RedisClient.FlushAll(r.Context()).Err()
	if err != nil {
		h.App.ErrorLog.Println(err)
		msg = "Could not empty cache"
	}

	h.render(w, r, views.CacheEmpty(msg, err))
}
