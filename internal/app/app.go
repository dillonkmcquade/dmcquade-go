package app

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type App struct {
	Views  *template.Template
	Router *chi.Mux
	Log    *log.Logger
	Data   AppData
}

// serve static assets like css, images, JS, etc.
func (h *App) ServeStatic() {
	fs := http.FileServer(http.Dir("web/static/"))
	h.Router.Handle("/static/*", http.StripPrefix("/static/", fs))
}

func (h *App) NotFound(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("404").Parse(`<h1 style="text-align: center;">404 Not found</h1>`)
	if err != nil {
		h.Log.Println(err)
	}
	rw.WriteHeader(404)
	err = tmpl.Execute(rw, nil)
	if err != nil {
		h.Log.Println(err)
	}
}

func (h *App) loadData() {
	// load data into struct
	b, err := os.ReadFile("internal/data/data.json")
	if err != nil {
		h.Log.Fatalf("Failed to read data from internal/data/data.json: %s ", err)
	}
	err = json.Unmarshal(b, &h.Data.Projects)
	if err != nil {
		h.Log.Fatalf("Failed to unmarshal JSON data from internal/data/data.json: %s ", err)
	}

	b, err = os.ReadFile("internal/data/projectDetails.json")
	if err != nil {
		h.Log.Fatalf("Failed to read data from internal/data/projectDetails.json: %s ", err)
	}
	err = json.Unmarshal(b, &h.Data.ProjectDetails)
	if err != nil {
		h.Log.Fatalf("Failed to unmarshal data from internal/data/projectDetails.json: %s ", err)
	}

	h.Views = template.Must(template.ParseGlob("template/*.html"))
}

/* Load data from json files and parse templates */
func (h *App) Bootstrap() {
	h.loadData()

	h.Log = log.New(os.Stdout, "", log.LstdFlags)

	h.Router = chi.NewRouter()

	h.middleware()
}

func NewApp() *App {
	return &App{
		Data: AppData{
			// ProjectDetails
			// Projects
			Skills: [10]string{
				"Go",
				"Typescript",
				"HTML",
				"CSS",
				"MongoDB",
				"Node.js",
				"Linux",
				"React",
				"Docker",
				"PostgreSQL",
			},
		},
	}
}
