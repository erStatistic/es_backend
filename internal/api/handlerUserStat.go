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

func (cfg *Config) UserStatCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("UserStatCtx")
		userStatId := chi.URLParam(r, "userStatId")
		UserStatID, err := strconv.Atoi(userStatId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		userStat, err := cfg.DB.GetUserStat(r.Context(), int32(UserStatID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "UserStat not found"
			} else {
				msg = "Failed to get UserStat"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserStatKey, &userStat)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) CreateUserStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating user stat")
	type parameters struct {
		UserID      int32 `json:"userId"`
		CharacterID int32 `json:"characterId"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdUserStat, err := cfg.DB.CreateUserStat(r.Context(), database.CreateUserStatParams{
		UserID:      params.UserID,
		CharacterID: params.CharacterID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user stat", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User stat created", createdUserStat)
}

func (cfg *Config) DeleteUserStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting user stat")
	ctx := r.Context()
	userStat, ok := ctx.Value(UserStatKey).(*database.UserStat)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User stat not found", nil)
		return
	}
	err := cfg.DB.DeleteUserStat(r.Context(), int32(userStat.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user stat", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User stat deleted", nil)
}

func (cfg *Config) PatchUserStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching user stat")
	ctx := r.Context()
	userStat, ok := ctx.Value(UserStatKey).(*database.UserStat)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User stat not found", nil)
		return
	}
	type parameters struct {
		UserID      int32  `json:"userId"`
		CharacterID int32  `json:"characterId"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.UserID != 0 {
		userStat.UserID = params.UserID
	}
	if params.CharacterID != 0 {
		userStat.CharacterID = params.CharacterID
	}
	err := cfg.DB.PatchUserStat(r.Context(), database.PatchUserStatParams{
		UserID:      userStat.UserID,
		CharacterID: userStat.CharacterID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch user stat", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User stat patched", userStat)
}

func (cfg *Config) GetUserStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting user stat")
	ctx := r.Context()
	userStat, ok := ctx.Value(UserStatKey).(*database.UserStat)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User stat not found", nil)
		return
	}
	respondWithJson(w, http.StatusOK, "User stat", userStat)
}

func (cfg *Config) ListUserStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing user stat")

	userStats, err := cfg.DB.ListUserStat(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get UserStatByUserId", err)
		return
	}

	if userStats == nil {
		userStats = []database.UserStat{}
	}
	respondWithJson(w, http.StatusOK, "User stat", userStats)
}
