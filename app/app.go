package app

import (
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	server  *http.Server
	sigChan chan os.Signal
	done    chan bool
}

func New(data_template embed.FS, static embed.FS) *App {
	tmpl := templates(data_template)
	data := loadData(data_template)

	router := Router(tmpl, static, data)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return &App{
		server:  server,
		sigChan: make(chan os.Signal),
		done:    make(chan bool),
	}
}

// Default 404 not found response
func notFound(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("404").Parse(`<h1 style="text-align: center;">404 Not found</h1>`)
	if err != nil {
		log.Println(err)
	}
	rw.WriteHeader(404)
	err = tmpl.Execute(rw, nil)
	if err != nil {
		log.Println(err)
	}
}

func loadData(embedded_fs embed.FS) *AppData {
	appData := &AppData{
		Skills: []string{
			"Go",
			"Typescript",
			"HTML",
			"CSS",
			"Bash",
			"MongoDB",
			"Node.js",
			"Linux",
			"React",
			"Docker",
			"PostgreSQL",
			"Ansible",
		},
	}
	b, err := embedded_fs.ReadFile("data/data.json")
	if err != nil {
		log.Fatalf("Failed to read data from internal/data/data.json: %s ", err)
	}
	err = json.Unmarshal(b, &appData.Projects)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON data from internal/data/data.json: %s ", err)
	}

	return appData
}

// Parse templates and json data
func templates(embedded_fs embed.FS) *template.Template {
	return template.Must(template.ParseFS(embedded_fs, "template/*.html"))
}

func (a *App) Run() {
	// Shutdown when signal received
	go a.waitForShutdown()

	a.start()

	if success := <-a.done; success {
		log.Println("server shutdown successfully")
	}
}

func (a *App) start() {
	log.Printf("Listening on port %s", a.server.Addr)

	err := a.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %s", err)
	}
}

func (a *App) waitForShutdown() {
	// Listen for interrupt or terminate signals

	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("received %s, commencing graceful shutdown", <-a.sigChan)

	tc, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(tc); err != nil {
		log.Printf("on server shutdown: %s", err)
		a.done <- false
		return
	}

	a.done <- true
}
