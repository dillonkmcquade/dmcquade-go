package controllers

import (
	"net/http"
	"strconv"

	"github.com/dillonkmcquade/dmcquade-go/internal/app"
	"github.com/go-chi/chi/v5"
)

func ProjectsPage(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(rw, "Missing url parameters", http.StatusBadRequest)
			a.Log.Println("Missing url parameters")
			return
		}

		idx, err := strconv.Atoi(id)
		if err != nil {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			a.Log.Println(err)
			return
		}

		err = a.Views.ExecuteTemplate(rw, "project_page.html", a.Data.ProjectDetails[idx])
		if err != nil {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			a.Log.Println(err)
			return
		}
	}
}
