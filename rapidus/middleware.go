package rapidus

import "net/http"

func (r *Rapidus) SessionLoad(next http.Handler) http.Handler {
	r.InfoLog.Println("SessionLoad")
	return r.Session.LoadAndSave(next)
}
