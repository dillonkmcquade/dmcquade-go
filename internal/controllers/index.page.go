package controllers

import (
	"net/http"

	"github.com/dillonkmcquade/dmcquade-go/internal/app"
)

func Index(app *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := app.Views.ExecuteTemplate(rw, "home.html", app.Data)
		if err != nil {
			app.Log.Println(err)
		}
	}
}
