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

func (cfg *Config) TierCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("TierCtx")
		id := chi.URLParam(r, "tierId")
		TierID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Tier, err := cfg.DB.GetTier(r.Context(), int32(TierID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Tier not found"
			} else {
				msg = "Failed to get tier"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), tierKey, Tier)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListTiers(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing tiers")
	tiers, err := cfg.DB.ListTiers(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list tiers(ListTiers Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListTiers", err)
		return
	}

	if tiers == nil {
		tiers = []database.Tier{}
	}

	respondWithJson(w, http.StatusOK, "Tiers retrieved", tiers)
}
func (cfg *Config) GetTier(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting tier")
	ctx := r.Context()
	tier, ok := ctx.Value(tierKey).(*database.Tier)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Tier not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Tier retrieved", tier)
}

func (cfg *Config) CreateTier(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating tier")

	type parameters struct {
		ImageUrl string `json:"imageUrl"`
		Name     string `json:"name"`
		Mmr      int32  `json:"mmr"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdTier, err := cfg.DB.CreateTier(r.Context(), database.CreateTierParams{
		ImageUrl: params.ImageUrl,
		Name:     params.Name,
		Mmr:      params.Mmr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create tier", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Tier created", createdTier)
}

func (cfg *Config) DeleteTier(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting tier")
	ctx := r.Context()
	tier, ok := ctx.Value(tierKey).(*database.Tier)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Tier not found", nil)
		return
	}

	err := cfg.DB.DeleteTier(r.Context(), int32(tier.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete tier", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Tier deleted", nil)
}

func (cfg *Config) PatchTier(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching tier")
	ctx := r.Context()
	tier, ok := ctx.Value(tierKey).(*database.Tier)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Tier not found", nil)
		return
	}
	type parameters struct {
		ImageUrl string `json:"imageUrl"`
		Name     string `json:"name"`
		Mmr      int32  `json:"mmr"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.ImageUrl != "" {
		tier.ImageUrl = params.ImageUrl
	}
	if params.Name != "" {
		tier.Name = params.Name
	}
	if params.Mmr != 0 {
		tier.Mmr = params.Mmr
	}

	err := cfg.DB.PatchTier(r.Context(), database.PatchTierParams{
		ID:       tier.ID,
		ImageUrl: tier.ImageUrl,
		Name:     tier.Name,
		Mmr:      tier.Mmr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch tier", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Tier patched", tier)
}
