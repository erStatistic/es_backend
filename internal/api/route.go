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
				r.Get("/", cfg.GetCharacter)
				r.Patch("/", cfg.PatchCharacter)
				r.Delete("/", cfg.DeleteCharacter)
			})
		})
		r.Route("/weapons", func(r chi.Router) {
			r.Get("/", cfg.ListWeapons)
			r.Post("/", cfg.CreateWeapon)

			r.Route("/{code}", func(r chi.Router) {
				r.Use(cfg.WeaponCtx)
				r.Get("/", cfg.GetWeapon)
				r.Patch("/", cfg.PatchWeapon)
				r.Delete("/", cfg.DeleteWeapon)
			})
		})

		r.Route("/positions", func(r chi.Router) {
			r.Get("/", cfg.ListPositions)
			r.Post("/", cfg.CreatePosition)

			r.Route("/{positionId}", func(r chi.Router) {
				r.Use(cfg.PositionCtx)
				r.Get("/", cfg.GetPosition)
				r.Patch("/", cfg.PatchPosition)
				r.Delete("/", cfg.DeletePosition)
			})
		})
		r.Route("/clusters", func(r chi.Router) {
			r.Get("/", cfg.ListClusters)
			r.Post("/", cfg.CreateCluster)
			r.Route("/{clusterId}", func(r chi.Router) {
				r.Use(cfg.ClusterCtx)
				r.Get("/", cfg.GetCluster)
				r.Patch("/", cfg.PatchCluster)
				r.Delete("/", cfg.DeleteCluster)
			})
		})
		r.Route("/tiers", func(r chi.Router) {
			r.Get("/", cfg.ListTiers)
			r.Post("/", cfg.CreateTier)
			r.Route("/{tierId}", func(r chi.Router) {
				r.Use(cfg.TierCtx)
				r.Get("/", cfg.GetTier)
				r.Patch("/", cfg.PatchTier)
				r.Delete("/", cfg.DeleteTier)
			})
		})
		r.Route("/times", func(r chi.Router) {
			r.Get("/", cfg.ListTimes)
			r.Post("/", cfg.CreateTime)
			r.Route("/{timeId}", func(r chi.Router) {
				r.Use(cfg.TimesCtx)
				r.Get("/", cfg.GetTime)
				r.Patch("/", cfg.PatchTime)
				r.Delete("/", cfg.DeleteTime)
			})
		})
		r.Route("/cws", func(r chi.Router) {
			r.Get("/", cfg.ListCharacterWeapons)
			r.Post("/", cfg.CreateCharacterWeapon)
			r.Get("/stats", cfg.ListCharacterWeaponStats)
			r.Route("/{cwId}", func(r chi.Router) {
				r.Use(cfg.CharacterWeaponCtx)
				r.Get("/", cfg.GetCharacterWeapon)
				r.Patch("/", cfg.PatchCharacterWeapon)
				r.Delete("/", cfg.DeleteCharacterWeapon)
				r.Route("/stats", func(r chi.Router) {
					r.Use(cfg.CharacterWeaponStatCtx)
					r.Post("/", cfg.CreateCharacterWeaponStat)
					r.Get("/", cfg.GetCharacterWeaponStat)
					r.Patch("/", cfg.PatchCharacterWeaponStat)
					r.Delete("/", cfg.DeleteCharacterWeaponStat)
				})
			})
		})
		r.Route("/games", func(r chi.Router) {
			r.Get("/", cfg.ListGames)
			r.Post("/", cfg.CreateGame)
			r.Route("/{rank}", func(r chi.Router) {
				r.Get("/", cfg.GetListGameRank)
			})
			r.Route("/{gameCode}", func(r chi.Router) {
				r.Use(cfg.GameCtx)
				r.Get("/", cfg.GetGame)
				r.Patch("/", cfg.PatchGame)
				r.Delete("/", cfg.DeleteGame)
				r.Route("/teams", func(r chi.Router) {
					r.Get("/", cfg.ListGameTeams)
					r.Post("/", cfg.CreateGameTeam)
					r.Route("/{teamId}", func(r chi.Router) {
						r.Use(cfg.GameTeamCtx)
						r.Get("/", cfg.GetGameTeam)
						r.Patch("/", cfg.PatchGameTeam)
						r.Delete("/", cfg.DeleteGameTeam)
						r.Route("/cws", func(r chi.Router) {
							r.Get("/", cfg.ListGameSameTeamCWs)
						})
					})
				})
			})
		})
		r.Route("/gameTeamCws", func(r chi.Router) {
			r.Get("/", cfg.ListGameTeamCWs)
			r.Post("/", cfg.CreateGameTeamCW)
			r.Route("/{gtcwId}", func(r chi.Router) {
				r.Use(cfg.GameTeamCWCtx)
				r.Get("/", cfg.GetGameTeamCW)
				r.Patch("/", cfg.PatchGameTeamCW)
				r.Delete("/", cfg.DeleteGameTeamCW)
			})
		})
	})

	return r
}
