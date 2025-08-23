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

func (cfg *Config) PositionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("PositionCtx")
		id := chi.URLParam(r, "positionId")
		PositionID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Position, err := cfg.DB.GetPosition(r.Context(), int32(PositionID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Position not found"
			} else {
				msg = "Failed to get position"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), positionKey, &Position)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListPositions(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing positions")
	positions, err := cfg.DB.ListPositions(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list positions(ListPositions Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListPositions", err)
		return
	}

	if positions == nil {
		positions = []database.Position{}
	}

	respondWithJson(w, http.StatusOK, "Positions retrieved", positions)
}
func (cfg *Config) GetPosition(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting position")
	ctx := r.Context()
	position, ok := ctx.Value(positionKey).(*database.Position)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Position not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Position retrieved", position)
}

func (cfg *Config) CreatePosition(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating position")

	type parameters struct {
		Name     string `json:"name"`
		ImageUrl string `json:"imageUrl"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdPosition, err := cfg.DB.CreatePosition(r.Context(), database.CreatePositionParams{
		Name:     params.Name,
		ImageUrl: params.ImageUrl,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create position", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Position created", createdPosition)
}

func (cfg *Config) DeletePosition(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting position")
	ctx := r.Context()
	position, ok := ctx.Value(positionKey).(*database.Position)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Position not found", nil)
		return
	}

	err := cfg.DB.DeletePosition(r.Context(), int32(position.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete position", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Position deleted", nil)
}

func (cfg *Config) PatchPosition(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching position")
	ctx := r.Context()
	position, ok := ctx.Value(positionKey).(*database.Position)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Position not found", nil)
		return
	}
	type parameters struct {
		Name     string `json:"name"`
		ImageUrl string `json:"imageUrl"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.Name != "" {
		position.Name = params.Name
	}
	if params.ImageUrl != "" {
		position.ImageUrl = params.ImageUrl
	}

	err := cfg.DB.PatchPosition(r.Context(), database.PatchPositionParams{
		ID:       position.ID,
		ImageUrl: position.ImageUrl,
		Name:     position.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch position", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Position patched", position)
}
