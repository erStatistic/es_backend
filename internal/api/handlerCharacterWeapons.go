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

func (cfg *Config) CharacterWeaponCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("CharacterWeaponCtx")
		id := chi.URLParam(r, "cwId")
		CharacterWeaponID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		CharacterWeapon, err := cfg.DB.GetCharacterWeapon(r.Context(), int32(CharacterWeaponID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "CharacterWeapon not found"
			} else {
				msg = "Failed to get character weapon"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), characterWeaponKey, CharacterWeapon)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListCharacterWeapons(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing character weapons")
	characterWeapons, err := cfg.DB.ListCharacterWeapons(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list character weapons(ListCharacterWeapons Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListCharacterWeapons", err)
		return
	}

	if characterWeapons == nil {
		characterWeapons = []database.CharacterWeapon{}
	}

	respondWithJson(w, http.StatusOK, "Character weapons retrieved", characterWeapons)
}
func (cfg *Config) GetCharacterWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting character weapon")
	ctx := r.Context()
	characterWeapon, ok := ctx.Value(characterWeaponKey).(*database.CharacterWeapon)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character weapon not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Character weapon retrieved", characterWeapon)
}

func (cfg *Config) CreateCharacterWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating character weapon")

	type parameters struct {
		CharacterID int32 `json:"characterId"`
		WeaponID    int32 `json:"weaponId"`
		PositionID  int32 `json:"positionId"`
		ClusterID   int32 `json:"clusterId"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdCharacterWeapon, err := cfg.DB.CreateCharacterWeapon(r.Context(), database.CreateCharacterWeaponParams{
		CharacterID: params.CharacterID,
		WeaponID:    params.WeaponID,
		PositionID:  params.PositionID,
		ClusterID:   params.ClusterID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create character weapon", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character weapon created", createdCharacterWeapon)
}

func (cfg *Config) DeleteCharacterWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting character weapon")
	ctx := r.Context()
	characterWeapon, ok := ctx.Value(characterWeaponKey).(*database.CharacterWeapon)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character weapon not found", nil)
		return
	}

	err := cfg.DB.DeleteCharacterWeapon(r.Context(), int32(characterWeapon.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete character weapon", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character weapon deleted", nil)
}

func (cfg *Config) PatchCharacterWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching character weapon")
	ctx := r.Context()
	characterWeapon, ok := ctx.Value(characterWeaponKey).(*database.CharacterWeapon)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character weapon not found", nil)
		return
	}
	type parameters struct {
		PositionID int32 `json:"positionId"`
		ClusterID  int32 `json:"clusterId"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.PositionID != 0 {
		characterWeapon.PositionID = params.PositionID
	}
	if params.ClusterID != 0 {
		characterWeapon.ClusterID = params.ClusterID
	}

	err := cfg.DB.PatchCharacterWeapon(r.Context(), database.PatchCharacterWeaponParams{
		CharacterID: characterWeapon.CharacterID,
		PositionID:  characterWeapon.PositionID,
		ClusterID:   characterWeapon.ClusterID,
		WeaponID:    characterWeapon.WeaponID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch character weapon", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character weapon patched", characterWeapon)
}
