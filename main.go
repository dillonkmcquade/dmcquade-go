package main

import (
	"embed"

	"github.com/dillonkmcquade/dmcquade-go/app"
)

//go:embed web/static
var static embed.FS

//go:embed data/* template/*
var data_template embed.FS

func main() {
	// Initialize the application
	app := app.New(data_template, static)
	app.Run()
}
