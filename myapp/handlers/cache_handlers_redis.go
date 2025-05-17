package handlers

import (
	"myapp/views"
	"net/http"
)

func (h *Handlers) CacheDemoRedis(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, views.CachePageRedis())
}

func (h *Handlers) CacheSaveRedis(w http.ResponseWriter, r *http.Request) {
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

	h.render(w, r, views.CacheSaveRedis(msg, err))
}

func (h *Handlers) CacheGetRedis(w http.ResponseWriter, r *http.Request) {
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

	h.render(w, r, views.CacheGetRedis(msg, err))
}

func (h *Handlers) CacheDeleteRedis(w http.ResponseWriter, r *http.Request) {
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

	h.render(w, r, views.CacheDeleteRedis(msg, err))
}

func (h *Handlers) CacheEmptyRedis(w http.ResponseWriter, r *http.Request) {
	msg := "Emptied cache"
	err := h.App.RedisClient.FlushAll(r.Context()).Err()
	if err != nil {
		h.App.ErrorLog.Println(err)
		msg = "Could not empty cache"
	}

	h.render(w, r, views.CacheEmptyRedis(msg, err))
}
