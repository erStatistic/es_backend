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

func (cfg *Config) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("UserCtx")
		userId := chi.URLParam(r, "userId")
		UserID, err := strconv.Atoi(userId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		user, err := cfg.DB.GetUserByUserNum(r.Context(), int32(UserID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "User not found"
			} else {
				msg = "Failed to get user"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListUsers(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing users")
	users, err := cfg.DB.ListUsers(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list users(ListUsers Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListUsers", err)
		return
	}
	if users == nil {
		users = []database.User{}
	}
	respondWithJson(w, http.StatusOK, "Users retrieved", users)
}

func (cfg *Config) GetUser(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting user")
	ctx := r.Context()
	user, ok := ctx.Value(UserKey).(*database.User)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "User retrieved", user)
}

func (cfg *Config) CreateUser(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating user")
	type parameters struct {
		Nickname string `json:"nickname"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}

	user, err := cfg.erapiClient.UserByNickname(params.Nickname)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user nickname", err)
		return
	}

	createdUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Nickname: params.Nickname,
		UserNum:  int32(user.UserNum),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	userstats, err := cfg.erapiClient.UserStatByUserId(int(user.UserNum))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user num", err)
		return
	}
	for _, userstat := range userstats {
		_, err := cfg.DB.CreateUserStat(r.Context(), database.CreateUserStatParams{
			UserID:      int32(createdUser.UserNum),
			CharacterID: int32(userstat.CharacterCode),
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create UserStat", err)
			return
		}
	}

	respondWithJson(w, http.StatusOK, "User created", createdUser)
}

func (cfg *Config) DeleteUser(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting user")
	ctx := r.Context()
	user, ok := ctx.Value(UserKey).(*database.User)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User not found", nil)
		return
	}
	err := cfg.DB.DeleteUser(r.Context(), int32(user.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User deleted", nil)
}

func (cfg *Config) PatchUser(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching user")
	ctx := r.Context()
	user, ok := ctx.Value(UserKey).(*database.User)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User not found", nil)
		return
	}
	type parameters struct {
		Nickname string `json:"nickname"`
		UserNum  int32  `json:"userNum"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.Nickname != "" {
		user.Nickname = params.Nickname
		u, err := cfg.erapiClient.UserByNickname(params.Nickname)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Invalid user nickname", err)
			return
		}
		user.UserNum = int32(u.UserNum)
	}

	err := cfg.DB.PatchUser(r.Context(), database.PatchUserParams{
		ID:       user.ID,
		Nickname: user.Nickname,
		UserNum:  user.UserNum,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch user", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User patched", user)
}

func (cfg *Config) ListUserTop3(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing user top3")
	ctx := r.Context()
	user, ok := ctx.Value(UserKey).(*database.User)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "User not found", nil)
		return
	}

	usertop3, err := cfg.DB.GetUserStatbyUserId(r.Context(), int32(user.UserNum))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get UserStatByUserId", err)
		return
	}
	respondWithJson(w, http.StatusOK, "User top3", usertop3)
}
