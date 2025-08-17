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

func (cfg *Config) CharacterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("CharacterCtx")
		code := chi.URLParam(r, "code")
		CharacterID, err := strconv.Atoi(code)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Character, err := cfg.DB.GetCharacter(r.Context(), int32(CharacterID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Character not found"
			} else {
				msg = "Failed to get character"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), characterKey, Character)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListCharacters(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing characters")
	characters, err := cfg.DB.ListCharacters(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list characters(ListCharacters Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListCharacaters", err)
		return
	}

	if characters == nil {
		characters = []database.Character{}
	}

	respondWithJson(w, http.StatusOK, "Characters retrieved", characters)
}
func (cfg *Config) GetCharacter(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting character")
	ctx := r.Context()
	character, ok := ctx.Value(characterKey).(*database.Character)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Character retrieved", character)
}

func (cfg *Config) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating character")

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
	createdCharacter, err := cfg.DB.CreateCharacter(r.Context(), database.CreateCharacterParams{
		Code:     params.Code,
		ImageUrl: params.ImageUrl,
		NameKr:   params.NameKr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create character", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character created", createdCharacter)
}

func (cfg *Config) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting character")
	ctx := r.Context()
	character, ok := ctx.Value(characterKey).(*database.Character)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character not found", nil)
		return
	}

	err := cfg.DB.DeleteCharacter(r.Context(), int32(character.Code))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete character", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character deleted", nil)
}

func (cfg *Config) PatchCharacter(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching character")
	ctx := r.Context()
	character, ok := ctx.Value(characterKey).(*database.Character)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Character not found", nil)
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
		character.NameKr = params.NameKr
	}
	if params.ImageUrl != "" {
		character.ImageUrl = params.ImageUrl
	}

	err := cfg.DB.PatchCharacter(r.Context(), database.PatchCharacterParams{
		Code:     character.Code,
		ImageUrl: character.ImageUrl,
		NameKr:   character.NameKr,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch character", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Character patched", character)
}
