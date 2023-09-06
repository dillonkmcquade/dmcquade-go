package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
)

type App struct {
	tmpl   *template.Template
	router *chi.Mux
	log    *log.Logger
	server *http.Server
}

type Home struct {
	Projects [4]Project
	Skills   [11]string
}

type Project struct {
	Src         string `json:"src"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Github      string `json:"github"`
	Youtube     string `json:"youtube"`
	Id          int    `json:"id"`
	NoInfo      bool   `json:"noInfo"`
}

type ProjectDetailSection struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ProjectDetails struct {
	Id       int                    `json:"id"`
	Title    string                 `json:"title"`
	Image    string                 `json:"image"`
	Sections []ProjectDetailSection `json:"sections"`
}

// load templates into struct

// serve static assets like css, images, JS, etc.
func (h *App) serveStatic() {
	fs := http.FileServer(http.Dir("web/static/"))
	h.router.Handle("/static/*", http.StripPrefix("/static/", fs))
}

// cleaner wrapper for Get
func (h *App) Get(path string, f http.HandlerFunc) {
	h.router.Get(path, f)
}

func (h *App) index(rw http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("internal/data/data.json")
	if err != nil {
		h.log.Println(err)
	}
	d := Home{
		Skills: [11]string{
			"Python",
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
	}
	err = json.Unmarshal(b, &d.Projects)
	if err != nil {
		h.log.Println(err)
	}
	h.tmpl = template.Must(template.ParseFiles("web/template/index.html"))
	err = h.tmpl.Execute(rw, d)
	if err != nil {
		h.log.Println(err)
	}
}

func (h *App) notFound(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("404").Parse(`<h1 style="text-align: center;">404 Not found</h1>`)
	if err != nil {
		h.log.Println(err)
	}
	rw.WriteHeader(404)
	tmpl.Execute(rw, nil)
}

func (h *App) projects(rw http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(rw, "Missing url parameters", http.StatusBadRequest)
		h.log.Println("Missing url parameters")
		return
	}

	b, err := os.ReadFile("internal/data/projectDetails.json")
	if err != nil {
		http.Error(rw, "Unable to retrieve project details", http.StatusNotFound)
		h.log.Println(err)
		return
	}

	d := [3]ProjectDetails{}

	err = json.Unmarshal(b, &d)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		h.log.Println(err)
		return
	}

	idx, err := strconv.Atoi(id)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		h.log.Println(err)
		return
	}

	h.tmpl = template.Must(template.ParseFiles("web/template/index.html", "web/template/project_page.html"))
	err = h.tmpl.Execute(rw, d[idx])
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		h.log.Println(err)
		return
	}
}

func main() {
	// create logger
	l := log.New(os.Stdout, "", log.LstdFlags)

	// Initialize the application
	app := App{
		log:    l,
		router: chi.NewRouter(),
		server: &http.Server{
			Addr:         ":8080",
			ErrorLog:     l,
			IdleTimeout:  120 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	}

	secureMiddleware := secure.New(secure.Options{
		ContentSecurityPolicy: "default-src 'self';base-uri 'self';font-src googleapis.com gstatic.com https: data:;form-action 'self';frame-ancestors 'self';img-src 'self' data:;object-src 'none';script-src 'self' https: unpkg.com 'unsafe-inline';script-src-attr 'none';style-src 'self'  https: 'unsafe-inline';upgrade-insecure-requests",
		ReferrerPolicy:        "no-referrer",
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
	})

	// middleware
	app.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	app.router.Use(secureMiddleware.Handler)
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.router.Use(middleware.Compress(5, "text/html", "text/css", "text/plain", "text/javascript", "image/vnd.microsoft.icon", "image/png", "image/jpeg"))
	app.serveStatic()

	// Index
	app.Get("/", app.index)

	// Project detail page
	app.Get("/project/{id}", app.projects)

	// Catch 404
	app.router.NotFound(app.notFound)

	// load router into custom server struct
	app.server.Handler = app.router

	// non blocking server
	go func() {
		app.log.Printf("Listening on port %s\n", app.server.Addr)

		err := app.server.ListenAndServe()
		if err != nil {
			app.log.Fatal(err)
		}
	}()

	// Notify on Interrupt/kill
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	app.log.Printf("Received %s, commencing graceful shutdown", <-sigChan)

	// graceful shutdown
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.server.Shutdown(tc); err != nil {
		app.log.Println(err)
	}
}
