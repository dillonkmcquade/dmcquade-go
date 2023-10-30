package app

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
)

func (h *App) middleware() {
	if h.Router == nil {
		panic("router not initialized")
	}
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
	h.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	h.Router.Use(middleware.Logger)
	h.Router.Use(middleware.Recoverer)
	h.Router.Use(secureMiddleware.Handler)
	h.Router.Use(middleware.Compress(5, "text/html", "text/css", "text/plain", "text/javascript", "image/vnd.microsoft.icon", "image/png", "image/jpeg"))
}
