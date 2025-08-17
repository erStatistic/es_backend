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
		r.Route("/weapons", func(r chi.Router) {
			r.Get("/", cfg.ListWeapons)
			r.Post("/", cfg.CreateWeapon)

			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.WeaponCtx)
				r.Get("/", cfg.GetWeapon) // GET /weapons/{code}
				r.Patch("/", cfg.PatchWeapon)
				r.Delete("/", cfg.DeleteWeapon)
			})
		})

		r.Route("/positions", func(r chi.Router) {
			r.Get("/", cfg.ListPositions)
			r.Post("/", cfg.CreatePosition)

			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.PositionCtx)
				r.Get("/", cfg.GetPosition) // GET /positions/{code}
				r.Patch("/", cfg.PatchPosition)
				r.Delete("/", cfg.DeletePosition)
			})
		})
		r.Route("/clusters", func(r chi.Router) {
			r.Get("/", cfg.ListClusters)
			r.Post("/", cfg.CreateCluster)
			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.ClusterCtx)
				r.Get("/", cfg.GetCluster) // GET /clusters/{code}
				r.Patch("/", cfg.PatchCluster)
				r.Delete("/", cfg.DeleteCluster)
			})
		})
		r.Route("/tiers", func(r chi.Router) {
			r.Get("/", cfg.ListTiers)
			r.Post("/", cfg.CreateTier)
			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.TierCtx)
				r.Get("/", cfg.GetTier) // GET /tiers/{code}
				r.Patch("/", cfg.PatchTier)
				r.Delete("/", cfg.DeleteTier)
			})
		})
		r.Route("/times", func(r chi.Router) {
			r.Get("/", cfg.ListTimes)
			r.Post("/", cfg.CreateTime)
			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.TimesCtx)
				r.Get("/", cfg.GetTime) // GET /times/{code}
				r.Patch("/", cfg.PatchTime)
				r.Delete("/", cfg.DeleteTime)
			})
		})
		r.Route("/cws", func(r chi.Router) {
			r.Get("/", cfg.ListCharacterWeapons)
			r.Post("/", cfg.CreateCharacterWeapon)
			r.Get("/stats", cfg.ListCharacterWeaponStats)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(cfg.CharacterWeaponCtx)
				r.Get("/", cfg.GetCharacterWeapon) // GET /character_weapons/{code}
				r.Patch("/", cfg.PatchCharacterWeapon)
				r.Delete("/", cfg.DeleteCharacterWeapon)
				r.Route("/stats", func(r chi.Router) {
					r.Get("/", cfg.GetCharacterWeaponStat)
					r.Patch("/", cfg.PatchCharacterWeaponStat)
					r.Delete("/", cfg.DeleteCharacterWeaponStat)
				})
			})
		})
	})

	return r
}
