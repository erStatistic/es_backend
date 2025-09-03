package rumiapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kaeba0616/es_backend/internal/database"
)

func (cfg *Config) GameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("GameCtx")
		id := chi.URLParam(r, "gameCode")
		GameID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Game, err := cfg.DB.GetGame(r.Context(), int64(GameID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Game not found"
			} else {
				msg = "Failed to get game"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), GameKey, &Game)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) CreateGame(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating game")
	type parameters struct {
		GameCode   int64              `json:"game_code"`
		AverageMmr int32              `json:"average_mmr"`
		StartedAt  pgtype.Timestamptz `json:"started_at"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}

	createdGame, err := cfg.DB.CreateGame(r.Context(), database.CreateGameParams{
		GameCode:   params.GameCode,
		AverageMmr: params.AverageMmr,
		StartedAt:  params.StartedAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create game", err)
		return
	}

	respondWithJson(w, http.StatusOK, "Game retrieved", toGameDTO(createdGame))
}

func (cfg *Config) GetGame(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting game")
	ctx := r.Context()
	game, ok := ctx.Value(GameKey).(*database.Game)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game not found", nil)
		return
	}
	respondWithJson(w, http.StatusOK, "Game retrieved", toGameDTO(*game))
}

func (cfg *Config) ListGames(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing games")
	games, err := cfg.DB.ListGames(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list games(ListGames Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListGames", err)
		return
	}
	if games == nil {
		games = []database.Game{}
	}
	respondWithJson(w, http.StatusOK, "Games retrieved", Map(games, toGameDTO))
}

func (cfg *Config) PatchGame(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching game")
	ctx := r.Context()
	game, ok := ctx.Value(GameKey).(*database.Game)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game not found", nil)
		return
	}
	type parameters struct {
		GameCode   int64     `json:"game_code"`
		AverageMmr int32     `json:"average_mmr"`
		StartedAt  time.Time `json:"started_at"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.GameCode != 0 {
		game.GameCode = params.GameCode
	}

	if params.AverageMmr != 0 {
		game.AverageMmr = params.AverageMmr
	}

	if !params.StartedAt.IsZero() {
		game.StartedAt.Time = params.StartedAt
	}

	err := cfg.DB.PatchGame(r.Context(), database.PatchGameParams{
		GameCode:   game.GameCode,
		AverageMmr: game.AverageMmr,
		StartedAt:  game.StartedAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch game", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game patched", toGameDTO(*game))
}

func (cfg *Config) DeleteGame(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting game")
	ctx := r.Context()
	game, ok := ctx.Value(GameKey).(*database.Game)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game not found", nil)
		return
	}

	err := cfg.DB.DeleteGame(r.Context(), game.GameCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete game", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game deleted", nil)
}

func (cfg *Config) TruncateGames(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Truncating games")
	err := cfg.DB.TruncateGames(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to truncate games", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Games truncated", nil)
}
