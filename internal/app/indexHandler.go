package app

import "net/http"

func index(app *App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := app.Views.ExecuteTemplate(rw, "home.html", app.Data)
		if err != nil {
			app.Log.Println(err)
		}
	}
}
