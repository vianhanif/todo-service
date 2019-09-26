package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vianhanif/todo-service/config"
	"github.com/vianhanif/todo-service/internal/app"
	"github.com/vianhanif/todo-service/internal/server/http/handlers"
	"github.com/vianhanif/todo-service/internal/server/http/middleware"
)

func addInstrument(r chi.Router, method, path string, handler http.HandlerFunc) {
	r.With(
		middleware.Instrument(method, path),
	).Method(method, path, handler)
}

func (hs *Server) compileRouter(config *config.Config) chi.Router {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Access-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)

	// A good base middleware stack
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	r.Use(middleware.Recoverer())
	r.Use(app.InjectorMiddleware(hs.app))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(chimw.Timeout(60 * time.Second))

	r.Handle("/metrics", promhttp.Handler())

	// PRIVATE endpoints
	r.Group(func(g chi.Router) {
		g.With(
			middleware.Private(config),
		).Route("/private/v1", func(r chi.Router) {
			addInstrument(r, "GET", "/todos", handlers.List())
			addInstrument(r, "POST", "/todos", handlers.Create())
			addInstrument(r, "GET", "/todos/:id", handlers.Find())
			// addInstrument(r, "PUT", "/todos/:id", handlers.Create())
			// addInstrument(r, "DELETE", "/todos/:id", handlers.Create())
		})
	})

	return r
}
