package rapidus

import "net/http"

func (r *Rapidus) SessionLoad(next http.Handler) http.Handler {
	return r.Session.LoadAndSave(next)
}
