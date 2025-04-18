package rapidus

import (
	"github.com/fouched/toolkit/v2"
	"net/http"
)

func (r *Rapidus) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	var t toolkit.Tools
	return t.WriteJSON(w, status, data, headers...)
}
