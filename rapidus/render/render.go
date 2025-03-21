package render

import (
	"github.com/a-h/templ"
	"net/http"
)

func Template(w http.ResponseWriter, r *http.Request, template templ.Component) error {

	return template.Render(r.Context(), w)
}
