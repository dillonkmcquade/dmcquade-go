package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dillonkmcquade/dmcquade-go/internal/app"
)

func main() {
	// Initialize the application
	app := app.NewApp()
	app.Bootstrap()
	app.ServeStatic()

	server := &http.Server{
		Addr:         ":8080",
		ErrorLog:     app.Log,
		Handler:      app.Router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// non blocking server
	go func() {
		app.Log.Printf("Listening on port %s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.Log.Printf("An error occurred while starting the server: %s", err)
		}
	}()

	// Notify on Interrupt/kill
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	app.Log.Printf("Received %s, commencing graceful shutdown", <-sigChan)

	// graceful shutdown
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
