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

func (cfg *Config) WeaponCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("WeaponCtx")
		code := chi.URLParam(r, "code")
		WeaponID, err := strconv.Atoi(code)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Weapon, err := cfg.DB.GetWeapon(context.Background(), int32(WeaponID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Weapon not found"
			} else {
				msg = "Failed to get weapon"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), weaponKey, Weapon)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListWeapons(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing weapons")
	weapons, err := cfg.DB.ListWeapons(context.Background())
	if err != nil {
		cfg.Log.Error("Failed to list weapons(ListWeapons Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListWeapons", err)
		return
	}

	if weapons == nil {
		weapons = []database.Weapon{}
	}

	respondWithJson(w, http.StatusOK, "Weapons retrieved", weapons)
}
func (cfg *Config) GetWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting weapon")
	ctx := r.Context()
	weapon, ok := ctx.Value(weaponKey).(database.Weapon)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Weapon not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Weapon retrieved", weapon)
}

func (cfg *Config) CreateWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating weapon")

	type parameters struct {
		Code     int32  `json:"code"`
		NameKr   string `json:"nameKr"`
		ImageUrl string `json:"imageUrl"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdWeapon, err := cfg.DB.CreateWeapon(context.Background(), database.CreateWeaponParams{
		Code:     params.Code,
		ImageUrl: params.ImageUrl,
		NameKr:   params.NameKr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create weapon", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Weapon created", createdWeapon)
}

func (cfg *Config) DeleteWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting weapon")
	ctx := r.Context()
	weapon, ok := ctx.Value(weaponKey).(database.Weapon)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Weapon not found", nil)
		return
	}

	err := cfg.DB.DeleteWeapon(context.Background(), int32(weapon.Code))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete weapon", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Weapon deleted", nil)
}

func (cfg *Config) PatchWeapon(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching weapon")
	ctx := r.Context()
	weapon, ok := ctx.Value(weaponKey).(database.Weapon)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Weapon not found", nil)
		return
	}
	type parameters struct {
		NameKr   string `json:"nameKr"`
		ImageUrl string `json:"imageUrl"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.NameKr != "" {
		weapon.NameKr = params.NameKr
	}
	if params.ImageUrl != "" {
		weapon.ImageUrl = params.ImageUrl
	}

	err := cfg.DB.PatchWeapon(context.Background(), database.PatchWeaponParams{
		Code:     weapon.Code,
		ImageUrl: weapon.ImageUrl,
		NameKr:   weapon.NameKr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch weapon", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Weapon patched", weapon)
}
