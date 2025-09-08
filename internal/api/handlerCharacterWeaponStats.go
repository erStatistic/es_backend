package rumiapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kaeba0616/es_backend/internal/database"
)

func (cfg *Config) CharacterWeaponStatCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("CharacterWeaponStatCtx")
		id := chi.URLParam(r, "cwId")
		CharacterWeaponID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		CharacterWeaponStat, err := cfg.DB.GetCharacterWeaponStat(r.Context(), int32(CharacterWeaponID))
		fmt.Println(CharacterWeaponStat)
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "CharacterWeaponStat not found"
			} else {
				msg = "Failed to get character weapon stat"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), CharacterWeaponStatKey, &CharacterWeaponStat)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) CreateCharacterWeaponStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating character weapon stat")
	type parameters struct {
		ID  int32 `json:"id"`
		Atk int32 `json:"atk"`
		Def int32 `json:"def"`
		Cc  int32 `json:"cc"`
		Spd int32 `json:"spd"`
		Sup int32 `json:"sup"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdCharacterWeaponStat, err := cfg.DB.CreateCharacterWeaponStat(r.Context(), database.CreateCharacterWeaponStatParams{
		CwID: params.ID,
		Atk:  params.Atk,
		Def:  params.Def,
		Cc:   params.Cc,
		Spd:  params.Spd,
		Sup:  params.Sup,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create character weapon stat", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character weapon stat created", createdCharacterWeaponStat)
}

func (cfg *Config) GetCharacterWeaponStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting character weapon stat")
	ctx := r.Context()
	characterWeaponStat, ok := ctx.Value(CharacterWeaponStatKey).(*database.CharacterWeaponStat)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character weapon stat not found", nil)
		return
	}
	respondWithJson(w, http.StatusOK, "Character weapon stat retrieved", characterWeaponStat)
}

func (cfg *Config) ListCharacterWeaponStats(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing character weapon stats")
	characterWeaponStats, err := cfg.DB.ListCharacterWeaponStats(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list character weapon stats(ListCharacterWeaponStats Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListCharacterWeaponStats", err)
		return
	}

	if characterWeaponStats == nil {
		characterWeaponStats = []database.CharacterWeaponStat{}
	}

	respondWithJson(w, http.StatusOK, "Character weapon stats retrieved", characterWeaponStats)
}

func (cfg *Config) PatchCharacterWeaponStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching character weapon stat")
	ctx := r.Context()
	characterWeaponStat, ok := ctx.Value(CharacterWeaponStatKey).(*database.CharacterWeaponStat)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character weapon stat not found", nil)
		return
	}
	type parameters struct {
		Atk int32 `json:"atk"`
		Def int32 `json:"def"`
		Cc  int32 `json:"cc"`
		Spd int32 `json:"spd"`
		Sup int32 `json:"sup"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.Atk != 0 {
		characterWeaponStat.Atk = params.Atk
	}
	if params.Def != 0 {
		characterWeaponStat.Def = params.Def
	}
	if params.Cc != 0 {
		characterWeaponStat.Cc = params.Cc
	}
	if params.Spd != 0 {
		characterWeaponStat.Spd = params.Spd
	}
	if params.Sup != 0 {
		characterWeaponStat.Sup = params.Sup
	}

	err := cfg.DB.PatchCharacterWeaponStat(r.Context(), database.PatchCharacterWeaponStatParams{
		CwID: characterWeaponStat.CwID,
		Atk:  characterWeaponStat.Atk,
		Def:  characterWeaponStat.Def,
		Cc:   characterWeaponStat.Cc,
		Spd:  characterWeaponStat.Spd,
		Sup:  characterWeaponStat.Sup,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch character weapon stat", err)
		return
	}

	respondWithJson(w, http.StatusOK, "Character weapon stat patched", characterWeaponStat)
}

func (cfg *Config) DeleteCharacterWeaponStat(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting character weapon stat")
	ctx := r.Context()
	characterWeaponStat, ok := ctx.Value(CharacterWeaponStatKey).(*database.CharacterWeaponStat)
	fmt.Println(characterWeaponStat)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character weapon stat not found", nil)
		return
	}

	err := cfg.DB.DeleteCharacterWeaponStat(r.Context(), characterWeaponStat.CwID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete character weapon stat", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character weapon stat deleted", nil)
}
