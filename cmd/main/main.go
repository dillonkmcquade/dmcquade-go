package main

import (
	"context"
	"encoding/json"
	"fmt"
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
func (h *App) parseTemplates(path string) {
	tmpl := template.Must(template.ParseFiles(path))
	h.tmpl = tmpl
}

// serve static assets like css, images, etc.
func (h *App) serveStatic() {
	fs := http.FileServer(http.Dir("web/static/"))
	h.router.Handle("/static/*", http.StripPrefix("/static/", fs))
}

// cleaner wrapper for Get
func (h *App) Get(path string, f http.HandlerFunc) {
	h.router.Get(path, f)
}

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)

	app := App{}
	app.log = l
	app.parseTemplates("web/views/layout.html")
	app.router = chi.NewRouter()
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.serveStatic()

	// Index
	app.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("internal/data/data.json")
		if err != nil {
			app.log.Println(err)
		}
		d := Home{
			Skills: []string{
				"JavaScript",
				"Typescript",
				"HTML",
				"CSS",
				"Express.js",
				"MongoDB",
				"Node.js",
				"styled-components",
				"linux",
				"React",
				"Docker",
				"Podman",
			},
		}
		err = json.Unmarshal(b, &d.Projects)
		if err != nil {
			app.log.Println(err)
		}
		err = app.tmpl.Execute(rw, d)
		if err != nil {
			fmt.Println(err)
		}
	})

	// Catch 404
	app.router.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("404").Parse(`<h1 style="text-align: center;">404 Not found</h1>`)
		if err != nil {
			fmt.Println(err)
		}
		rw.WriteHeader(404)
		tmpl.Execute(rw, "404")
	})

	// server opts
	server := http.Server{
		Addr:         ":3001",
		Handler:      app.router,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// non blocking server
	go func() {
		app.log.Printf("Listening on port %s\n", server.Addr)

		err := server.ListenAndServe()
		if err != nil {
			app.log.Fatal(err)
		}
	}()

	// Notify on Interrupt/kill
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	l.Printf("Received %s, commencing graceful shutdown", <-sigChan)

	// graceful shutdown
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		app.log.Println(err)
	}
}
