package routes

import (
	"net/http"

	"github.com/dillonkmcquade/portfolio-go/internal/views"
)

func Index(rw http.ResponseWriter, r *http.Request, t views.Templates) {
	t.Execute(rw, "layout", nil)
}
