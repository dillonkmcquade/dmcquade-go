package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	tmpl   *template.Template
	router *chi.Mux
	log    *log.Logger
	server *http.Server
}

type Home struct {
	Projects []Project
	Skills   []string
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
		Skills: []string{
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
			"Podman",
		},
	}
	err = json.Unmarshal(b, &d.Projects)
	if err != nil {
		h.log.Println(err)
	}
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

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)
	tmpl := template.Must(template.ParseFiles("web/views/layout.html"))
	app := App{
		tmpl:   tmpl,
		log:    l,
		router: chi.NewRouter(),
		server: &http.Server{
			Addr:         ":3001",
			ErrorLog:     l,
			IdleTimeout:  120 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	}
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.serveStatic()

	// Index
	app.Get("/", app.index)

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
