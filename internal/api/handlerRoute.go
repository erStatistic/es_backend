package rumiapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kaeba0616/es_backend/internal/database"
)

func (cfg *Config) UserRouteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("UserRouteCtx")
		routeId := chi.URLParam(r, "routeId")
		RouteID, err := strconv.Atoi(routeId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		route, err := cfg.DB.GetUserRoute(r.Context(), int32(RouteID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Route not found"
			} else {
				msg = "Failed to get Route"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), RouteKey, &route)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) GetUserRoute(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting user route")
	ctx := r.Context()
	userRoute, ok := ctx.Value(RouteKey).(*database.UserRoute)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User route not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "User route retrieved", userRoute)
}

func (cfg *Config) CreateUserRoute(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating user route")
	type parameters struct {
		RouteID     int32  `json:"routeId"`
		WeaponID    int32  `json:"weaponId"`
		CharacterID int32  `json:"characterId"`
		Title       string `json:"title"`
		Count       int32  `json:"count"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdUserRoute, err := cfg.DB.CreateUserRoute(r.Context(), database.CreateUserRouteParams{
		RouteID:     params.RouteID,
		WeaponID:    params.WeaponID,
		CharacterID: params.CharacterID,
		Title:       params.Title,
		Count:       params.Count,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user route", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User route created", createdUserRoute)
}

func (cfg *Config) DeleteUserRoute(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting user route")
	ctx := r.Context()
	userRoute, ok := ctx.Value(RouteKey).(*database.UserRoute)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User route not found", nil)
		return
	}

	err := cfg.DB.DeleteUserRoute(r.Context(), int32(userRoute.RouteID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user route", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User route deleted", nil)
}

func (cfg *Config) PatchUserRoute(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching user route")
	ctx := r.Context()
	userRoute, ok := ctx.Value(RouteKey).(*database.UserRoute)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User route not found", nil)
		return
	}
	type parameters struct {
		WeaponID    int32  `json:"weaponId"`
		CharacterID int32  `json:"characterId"`
		Title       string `json:"title"`
		Count       int32  `json:"count"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.WeaponID != 0 {
		userRoute.WeaponID = params.WeaponID
	}
	if params.CharacterID != 0 {
		userRoute.CharacterID = params.CharacterID
	}
	if params.Count != 0 {
		userRoute.Count = params.Count
	}
	if params.Title != "" {
		userRoute.Title = params.Title
	}

	err := cfg.DB.PatchUserRoute(r.Context(), database.PatchUserRouteParams{
		RouteID:     userRoute.RouteID,
		WeaponID:    userRoute.WeaponID,
		CharacterID: userRoute.CharacterID,
		Count:       userRoute.Count,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch user route", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User route patched", userRoute)
}

func (cfg *Config) ListUserRoutes(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing user routes")
	userRoutes, err := cfg.DB.ListUserRoutes(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list user routes(ListUserRoutes Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListUserRoutes", err)
		return
	}
	if userRoutes == nil {
		userRoutes = []database.UserRoute{}
	}
	respondWithJson(w, http.StatusOK, "User routes retrieved", userRoutes)
}
