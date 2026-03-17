package server

import (
	"io/fs"
	"net/http"
	"strings"
	"backstock/internal/handler"
	"backstock/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func bodySizeLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func New(s *store.Store, frontend fs.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(bodySizeLimit(1 << 20)) // 1 MB
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	items := handler.NewItemHandler(s)
	stock := handler.NewStockHandler(s)
	cats := handler.NewCategoryHandler(s)
	units := handler.NewUnitHandler(s)
	audit := handler.NewAuditHandler(s)
	locs := handler.NewLocationHandler(s)

	r.Route("/api", func(r chi.Router) {
		r.Get("/items", items.List)
		r.Post("/items", items.Create)
		r.Get("/items/{id}", items.Get)
		r.Put("/items/{id}", items.Update)
		r.Delete("/items/{id}", items.Delete)

		r.Get("/items/{id}/stock", stock.ListForItem)
		r.Post("/items/{id}/stock", stock.Add)

		r.Put("/stock/{stockID}", stock.Update)
		r.Delete("/stock/{stockID}", stock.Delete)
		r.Post("/stock/{stockID}/move", stock.Move)

		r.Get("/categories", cats.List)
		r.Post("/categories", cats.Create)
		r.Put("/categories/{id}", cats.Update)
		r.Delete("/categories/{id}", cats.Delete)

		r.Get("/units", units.List)
		r.Post("/units", units.Create)
		r.Put("/units/{id}", units.Update)
		r.Delete("/units/{id}", units.Delete)

		r.Get("/locations", locs.List)
		r.Post("/locations", locs.Create)
		r.Put("/locations/{id}", locs.Update)
		r.Delete("/locations/{id}", locs.Delete)

		r.Get("/audit", audit.List)
	})

	// SPA fallback
	fileServer := http.FileServer(http.FS(frontend))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Try to serve the file directly
		if path != "/" && !strings.HasPrefix(path, "/api") {
			f, err := frontend.Open(strings.TrimPrefix(path, "/"))
			if err == nil {
				f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		// Fallback to index.html for SPA routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	return r
}
