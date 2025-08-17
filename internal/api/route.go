package rumiapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (cfg *Config) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/characters", func(r chi.Router) {
			r.Get("/", cfg.ListCharacters)
			r.Post("/", cfg.CreateCharacter)

			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.CharacterCtx)
				r.Get("/", cfg.GetCharacter) // GET /characters/{code}
				r.Patch("/", cfg.PatchCharacter)
				r.Delete("/", cfg.DeleteCharacter)
			})
		})
	})

	return r
}
