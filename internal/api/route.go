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
		r.Route("/users", func(r chi.Router) {
			r.Get("/", cfg.ListUsers)
			r.Post("/", cfg.CreateUser)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(cfg.UserCtx)
				r.Get("/", cfg.GetUser)
				r.Patch("/", cfg.PatchUser)
				r.Delete("/", cfg.DeleteUser)
				r.Get("/top3", cfg.ListUserTop3)
			})
			r.Get("/search", cfg.GetUserByNickname)
			r.Route("/stats", func(r chi.Router) {
				r.Post("/", cfg.CreateUserStat)
				r.Get("/", cfg.ListUserStat)
				r.Route("/{userStatId}", func(r chi.Router) {
					r.Use(cfg.UserStatCtx)
					r.Get("/", cfg.GetUserStat)
					r.Patch("/", cfg.PatchUserStat)
					r.Delete("/", cfg.DeleteUserStat)
				})
			})
		})

		r.Route("/characters", func(r chi.Router) {
			r.Get("/", cfg.ListCharacters)
			r.Post("/", cfg.CreateCharacter)

			r.Route("/{characterId}", func(r chi.Router) {
				r.Use(cfg.CharacterCtx)
				r.Get("/", cfg.GetCharacter)
				r.Patch("/", cfg.PatchCharacter)
				r.Delete("/", cfg.DeleteCharacter)
				r.Get("/cws", cfg.ListCwsByCharacter)
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
			r.Get("/directory", cfg.ListCwDirectoryByCluster)
			r.Get("/by-cluster/{clusterId}", cfg.ListCwEntriesByCluster)
			r.Route("/{cwId}", func(r chi.Router) {
				r.Use(cfg.CharacterWeaponCtx)
				r.Get("/", cfg.GetCharacterWeapon)
				r.Patch("/", cfg.PatchCharacterWeapon)
				r.Delete("/", cfg.DeleteCharacterWeapon)
				r.Get("/overview", cfg.GetCwOverview)
			})
		})
		r.Route("/cw-stats", func(r chi.Router) {
			r.Get("/", cfg.ListCharacterWeaponStats)
			r.Post("/", cfg.CreateCharacterWeaponStat)
			r.Route("/{cwId}", func(r chi.Router) {
				r.Use(cfg.CharacterWeaponStatCtx)
				r.Get("/", cfg.GetCharacterWeaponStat)
				r.Patch("/", cfg.PatchCharacterWeaponStat)
				r.Delete("/", cfg.DeleteCharacterWeaponStat)
			})
		})
		r.Route("/games", func(r chi.Router) {
			r.Get("/", cfg.ListGames)
			r.Post("/", cfg.CreateGame)
			r.Delete("/", cfg.TruncateGames) // Truncate Games postman
			r.Route("/{gameCode}", func(r chi.Router) {
				r.Use(cfg.GameCtx)
				r.Get("/", cfg.GetGame)
				r.Patch("/", cfg.PatchGame)
				r.Delete("/", cfg.DeleteGame)
				r.Route("/teams", func(r chi.Router) {
					r.Route("/{teamId}", func(r chi.Router) {
						r.Get("/", cfg.GetGameTeam)
					})
				})
			})
		})
		r.Route("/gameTeams", func(r chi.Router) {
			r.Post("/", cfg.CreateGameTeam)
			r.Get("/", cfg.ListGameTeams)
			r.Delete("/", cfg.TruncateGameTeams) // Truncate GameTeams postman
			r.Route("/{gtId}", func(r chi.Router) {
				r.Use(cfg.GameTeamCtx)
				r.Get("/", cfg.GetGameTeamByID)
				r.Patch("/", cfg.PatchGameTeam)
				r.Delete("/", cfg.DeleteGameTeam)
				r.Route("/cws", func(r chi.Router) {
					r.Get("/", cfg.ListGameSameTeamCWs)
				})
			})
			r.Route("/ranks", func(r chi.Router) {
				r.Route("/{rank}", func(r chi.Router) {
					r.Use(cfg.GameRankCtx)
					r.Get("/", cfg.GetListGameTeamRank)
				})
			})
		})
		r.Route("/gameTeamCws", func(r chi.Router) {
			r.Get("/", cfg.ListGameTeamCWs)
			r.Post("/", cfg.CreateGameTeamCW)
			r.Delete("/", cfg.TruncateGameTeamCWs) // Truncate GameTeamcws postman
			r.Post("/multi", cfg.CreateGameTeamCWList)
			r.Route("/{gtcwId}", func(r chi.Router) {
				r.Use(cfg.GameTeamCWCtx)
				r.Get("/", cfg.GetGameTeamCW)
				r.Patch("/", cfg.PatchGameTeamCW)
				r.Delete("/", cfg.DeleteGameTeamCW)
			})
		})
		r.Route("/routes", func(r chi.Router) {
			r.Get("/", cfg.ListUserRoutes)
			r.Post("/", cfg.CreateUserRoute)
			r.Route("/{routeId}", func(r chi.Router) {
				r.Use(cfg.UserRouteCtx)
				r.Get("/", cfg.GetUserRoute)
				r.Patch("/", cfg.PatchUserRoute)
				r.Delete("/", cfg.DeleteUserRoute)
			})
		})
		r.Route("/analytics", func(r chi.Router) {
			r.Get("/combos/clusters", cfg.GetClusterCombos)
			r.Get("/cw/stats", cfg.GetCwStats)
			r.Get("/cw/top5", cfg.GetCwStatTop5)
			r.Get("/popular-comps", cfg.GetTopPopularComps)
			r.Post("/comp/metrics", cfg.GetCompMetrics)
			r.Post("/mv/refresh", cfg.RefreshMvTrioTeams)
			r.Route("/cw/{cwId}", func(r chi.Router) {
				r.Get("/stats", cfg.GetCwStatsByCw)
				r.Get("/trend", cfg.GetCwTrend)
				r.Get("/best-comps", cfg.GetBestCompsByCw)
			})
		})

	})

	return r
}
