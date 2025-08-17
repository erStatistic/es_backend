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

func (cfg *Config) TimesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("TimesCtx")
		id := chi.URLParam(r, "id")
		TimeID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Time, err := cfg.DB.GetTime(r.Context(), int32(TimeID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Time not found"
			} else {
				msg = "Failed to get time"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), timeKey, Time)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListTimes(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing times")
	times, err := cfg.DB.ListTimes(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list times(ListTimes Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListTimes", err)
		return
	}

	if times == nil {
		times = []database.Time{}
	}

	respondWithJson(w, http.StatusOK, "Times retrieved", times)
}
func (cfg *Config) GetTime(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting time")
	ctx := r.Context()
	time, ok := ctx.Value(timeKey).(*database.Time)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Time not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Time retrieved", time)
}

func (cfg *Config) CreateTime(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating time")

	type parameters struct {
		No        int32  `json:"no"`
		Name      string `json:"name"`
		Seconds   int32  `json:"seconds"`
		StartTime int32  `json:"startTime"`
		EndTime   int32  `json:"endTime"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdTime, err := cfg.DB.CreateTime(r.Context(), database.CreateTimeParams{
		No:        params.No,
		Name:      params.Name,
		Seconds:   params.Seconds,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create time", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Time created", createdTime)
}

func (cfg *Config) DeleteTime(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting time")
	ctx := r.Context()
	time, ok := ctx.Value(timeKey).(*database.Time)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Time not found", nil)
		return
	}

	err := cfg.DB.DeleteTime(r.Context(), int32(time.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete time", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Time deleted", nil)
}

func (cfg *Config) PatchTime(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching time")
	ctx := r.Context()
	time, ok := ctx.Value(timeKey).(*database.Time)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Time not found", nil)
		return
	}
	type parameters struct {
		Name      string `json:"name"`
		Seconds   int32  `json:"seconds"`
		StartTime int32  `json:"startTime"`
		EndTime   int32  `json:"endTime"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode requset Body", err)
		return
	}

	if params.Name != "" {
		time.Name = params.Name
	}
	if params.Seconds != 0 {
		time.Seconds = params.Seconds
	}
	if params.StartTime != 0 {
		time.StartTime = params.StartTime
	}
	if params.EndTime != 0 {
		time.EndTime = params.EndTime
	}

	err := cfg.DB.PatchTime(r.Context(), database.PatchTimeParams{
		ID:        time.ID,
		Name:      time.Name,
		Seconds:   time.Seconds,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch time", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Time patched", time)
}
