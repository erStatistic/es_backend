package rumiapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kaeba0616/es_backend/internal/database"
)

func (cfg *Config) GameTeamCWCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("GameTeamCWCtx")

		id := chi.URLParam(r, "gtcwId")
		gtcwID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}

		gameteamCw, err := cfg.DB.GetGameTeamCW(r.Context(), int32(gtcwID))
		if err != nil {
			respondWithError(w, http.StatusNotFound, "GameTeamCW not found", err)
			return
		}
		ctx := context.WithValue(r.Context(), GameTeamCWKey, &gameteamCw)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) CreateGameTeamCW(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating GameTeamCW")
	type parameters struct {
		CwID       int32 `json:"cw_id"`
		GameTeamID int32 `json:"gameteam_id"`
		Mmr        int32 `json:"mmr"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}

	createdGameTeamCW, err := cfg.DB.CreateGameTeamCW(r.Context(), database.CreateGameTeamCWParams{
		GameTeamID: params.GameTeamID,
		CwID:       params.CwID,
		Mmr:        params.Mmr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create GameTeamCW", err)
		return
	}

	respondWithJson(w, http.StatusOK, "GameTeamCW retrieved", createdGameTeamCW)
}

func (cfg *Config) ListGameTeamCWs(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting list game team cws")
	gameTeamCWs, err := cfg.DB.ListGameTeamCWs(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list game team cws(ListGameTeamCWs Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListGameTeamCWs", err)
		return
	}
	if gameTeamCWs == nil {
		gameTeamCWs = []database.GameTeamCw{}
	}
	respondWithJson(w, http.StatusOK, "Game team cws retrieved", gameTeamCWs)
}

func (cfg *Config) ListGameSameTeamCWs(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting list game same team cws")
	ctx := r.Context()
	gameteam, ok := ctx.Value(GameTeamKey).(*database.GameTeam)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "GameTeamCW not found", nil)
		return
	}
	gameTeamCWs, err := cfg.DB.ListGameSameTeamCWs(r.Context(), gameteam.ID)
	if err != nil {
		cfg.Log.Error("Failed to list game team cws(ListGameSameTeamCWs Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListGameSameTeamCWs", err)
		return
	}
	if gameTeamCWs == nil {
		gameTeamCWs = []database.GameTeamCw{}
	}
	respondWithJson(w, http.StatusOK, "Game team cws retrieved", gameTeamCWs)
}

func (cfg *Config) GetGameTeamCW(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting GameTeamCW")
	ctx := r.Context()
	gtcw, ok := ctx.Value(GameTeamCWKey).(*database.GameTeamCw)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "GameTeamCW not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "GameTeamCW retrieved", gtcw)
}

func (cfg *Config) PatchGameTeamCW(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching GameTeamCW")
	ctx := r.Context()
	gtcw, ok := ctx.Value(GameTeamCWKey).(*database.GameTeamCw)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "GameTeamCW not found", nil)
		return
	}

	type parameters struct {
		GameTeamID int32 `json:"gameteam_id"`
		CwID       int32 `json:"cw_id"`
		Mmr        int32 `json:"mmr"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	if params.GameTeamID != 0 {
		gtcw.GameTeamID = params.GameTeamID
	}
	if params.CwID != 0 {
		gtcw.CwID = params.CwID
	}

	err := cfg.DB.PatchGameTeamCW(r.Context(), database.PatchGameTeamCWParams{
		GameTeamID: gtcw.GameTeamID,
		CwID:       gtcw.CwID,
		Mmr:        params.Mmr,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch GameTeamCW", err)
		return
	}

	respondWithJson(w, http.StatusOK, "GameTeamCW retrieved", gtcw)
}

func (cfg *Config) DeleteGameTeamCW(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting GameTeamCW")
	ctx := r.Context()
	gtcw, ok := ctx.Value(GameTeamCWKey).(*database.GameTeamCw)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "GameTeamCW not found", nil)
		return
	}

	err := cfg.DB.DeleteGameTeamCW(r.Context(), gtcw.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete GameTeamCW", err)
		return
	}

	respondWithJson(w, http.StatusOK, "GameTeamCW deleted", nil)
}

func (cfg *Config) TruncateGameTeamCWs(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Truncating game team cws")
	err := cfg.DB.TruncateGameTeamCWs(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to truncate game team cws", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game team cws truncated", nil)
}
