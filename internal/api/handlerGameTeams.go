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

func (cfg *Config) GameTeamCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("GameTeamCtx")
		id := chi.URLParam(r, "gtId")
		GameTeamID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		GameTeam, err := cfg.DB.GetGameTeamByID(r.Context(), int32(GameTeamID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "GameTeam not found"
			} else {
				msg = "Failed to get game team"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), GameTeamKey, &GameTeam)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) CreateGameTeam(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating game team")
	type parameters struct {
		GameID         int64 `json:"gameId"`
		TeamID         int32 `json:"teamId"`
		GameRank       int32 `json:"gameRank"`
		TeamKills      int32 `json:"teamKills"`
		MonsterCredits int32 `json:"monsterCredits"`
		GainedMmr      int32 `json:"gainedMmr"`
		TeamAvgMmr     int32 `json:"teamAvgMmr"`
		TotalTime      int32 `json:"totalTime"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {

		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdGameTeam, err := cfg.DB.CreateGameTeam(r.Context(), database.CreateGameTeamParams{
		GameID:         params.GameID,
		TeamID:         params.TeamID,
		GameRank:       params.GameRank,
		TeamKills:      params.TeamKills,
		MonsterCredits: params.MonsterCredits,
		GainedMmr:      params.GainedMmr,
		TeamAvgMmr:     params.TeamAvgMmr,
		TotalTime:      params.TotalTime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create game team", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game team created", createdGameTeam)
}

func (cfg *Config) GetGameTeam(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting game team")

	id := chi.URLParam(r, "teamId")
	TeamID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
		return
	}
	GameID, err := strconv.Atoi(chi.URLParam(r, "gameId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
		return
	}
	GameTeam, err := cfg.DB.GetGameTeam(r.Context(), database.GetGameTeamParams{
		GameID: int64(GameID),
		TeamID: int32(TeamID),
	})
	if err != nil {
		var msg string
		if err == sql.ErrNoRows {
			msg = "GameTeam not found"
		} else {
			msg = "Failed to get game team"
		}
		respondWithError(w, http.StatusNotFound, msg, err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game team retrieved", toGameTeamDTO(GameTeam))
}

func (cfg *Config) ListGameTeams(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing game teams")
	gameTeams, err := cfg.DB.ListGameTeams(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list game teams(ListGameTeams Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListGameTeams", err)
		return
	}
	if gameTeams == nil {
		gameTeams = []database.GameTeam{}
	}
	respondWithJson(w, http.StatusOK, "Game teams retrieved", Map(gameTeams, toGameTeamDTO))
}

func (cfg *Config) GetGameTeamByID(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting game team by id")
	ctx := r.Context()
	gameTeam, ok := ctx.Value(GameTeamKey).(*database.GameTeam)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game team not found", nil)
		return
	}
	respondWithJson(w, http.StatusOK, "Game team retrieved", toGameTeamDTO(*gameTeam))
}

func (cfg *Config) PatchGameTeam(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching game team")
	ctx := r.Context()
	gameTeam, ok := ctx.Value(GameTeamKey).(*database.GameTeam)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game team not found", nil)
		return
	}
	type parameters struct {
		GameID         int64 `json:"gameId"`
		TeamID         int32 `json:"teamId"`
		GameRank       int32 `json:"gameRank"`
		TeamKills      int32 `json:"teamKills"`
		MonsterCredits int32 `json:"monsterCredits"`
		GainedMmr      int32 `json:"gainedMmr"`
		TeamAvgMmr     int32 `json:"teamAvgMmr"`
		TotalTime      int32 `json:"totalTime"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.GameID != 0 {
		gameTeam.GameID = params.GameID
	}
	if params.TeamID != 0 {
		gameTeam.TeamID = params.TeamID
	}
	if params.GameRank != 0 {
		gameTeam.GameRank = params.GameRank
	}
	if params.TeamKills != 0 {
		gameTeam.TeamKills = params.TeamKills
	}
	if params.MonsterCredits != 0 {
		gameTeam.MonsterCredits = params.MonsterCredits
	}
	if params.GainedMmr != 0 {
		gameTeam.GainedMmr = params.GainedMmr
	}
	if params.TeamAvgMmr != 0 {
		gameTeam.TeamAvgMmr = params.TeamAvgMmr
	}
	if params.TotalTime != 0 {
		gameTeam.TotalTime = params.TotalTime
	}

	err := cfg.DB.PatchGameTeam(r.Context(), database.PatchGameTeamParams{
		ID:             gameTeam.ID,
		GameID:         gameTeam.GameID,
		TeamID:         gameTeam.TeamID,
		GameRank:       gameTeam.GameRank,
		TeamKills:      gameTeam.TeamKills,
		MonsterCredits: gameTeam.MonsterCredits,
		GainedMmr:      gameTeam.GainedMmr,
		TeamAvgMmr:     gameTeam.TeamAvgMmr,
		TotalTime:      gameTeam.TotalTime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch game team", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game team patched", toGameTeamDTO(*gameTeam))
}

func (cfg *Config) DeleteGameTeam(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting game team")
	ctx := r.Context()
	gameTeam, ok := ctx.Value(GameTeamKey).(*database.GameTeam)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game team not found", nil)
		return
	}

	err := cfg.DB.DeleteGameTeam(ctx, int32(gameTeam.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete game team", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game team deleted", nil)
}

func (cfg *Config) GameRankCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("GameRankCtx")
		rank, err := strconv.Atoi(chi.URLParam(r, "rank"))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		GameTeams, err := cfg.DB.GetListGameTeamsByGameRank(r.Context(), int32(rank))
		if err != nil {
			cfg.Log.Error("Failed to list game teams(ListGameTeamsByGameRank Query)", "error", err)
			respondWithError(w, http.StatusInternalServerError, "DB error ListGameTeamsByGameRank", err)
			return
		}
		ctx := context.WithValue(r.Context(), GameTeamKey, &GameTeams)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) GetListGameTeamRank(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting list game rank")
	ctx := r.Context()
	gameTeams, ok := ctx.Value(GameTeamKey).([]database.GameTeam)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Game team not found", nil)
		return
	}
	respondWithJson(w, http.StatusOK, "Game teams by game rank retrieved", Map(gameTeams, toGameTeamDTO))
}

func (cfg *Config) TruncateGameTeams(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Truncating game teams")
	err := cfg.DB.TruncateGameTeams(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to truncate game teams", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Game teams truncated", nil)
}
