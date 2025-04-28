package render

import (
	"context"
	"github.com/a-h/templ"
	"github.com/justinas/nosurf"
	"net/http"
)

func Template(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	
	// Create a context and set value(s) that will be available to all templates
	ctx := context.WithValue(context.Background(), "CSRFToken", nosurf.Token(r))

	return template.Render(ctx, w)
}
