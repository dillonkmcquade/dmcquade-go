package app

import (
	"html/template"
	"log"
	"net/http"
)

func index(tmpl *template.Template, data *AppData) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(rw, "home.html", data)
		if err != nil {
			log.Println(err)
		}
	}
}
