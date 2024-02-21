package app

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
)

func Router(tmpl *template.Template, static embed.FS, data *AppData) *chi.Mux {
	app := chi.NewRouter()

	secureMiddleware := secure.New(secure.Options{
		ContentSecurityPolicy: `
            default-src 'self' data:;
            base-uri 'self';
            font-src googleapis.com gstatic.com https: data:;
            form-action 'self';frame-ancestors 'self';
            img-src 'self' data:;object-src 'none';
            script-src 'strict-dynamic' 'nonce-dmcquade-go' 'sha384-xcuj3WpfgjlKF+FXhSQFQ0ZNr39ln+hwjN3npfM9VBnUskLolQAcN80McRIVOPuO';
            script-src-attr 'none';
            style-src 'self' https: 'unsafe-inline';
            upgrade-insecure-requests`,
		ReferrerPolicy:       "no-referrer",
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	})

	// middleware
	app.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(secureMiddleware.Handler)
	app.Use(middleware.Compress(5, "text/html", "text/css", "text/plain", "text/javascript", "image/vnd.microsoft.icon", "image/png", "image/jpeg"))

	// Serve static files
	sub, err := fs.Sub(static, "web/static")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(sub))
	app.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Routes
	app.Get("/", index(tmpl, data))

	// 404
	app.NotFound(notFound)

	return app
}
