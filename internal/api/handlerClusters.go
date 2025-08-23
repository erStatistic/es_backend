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

func (cfg *Config) ClusterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Log.Info("ClusterCtx")
		id := chi.URLParam(r, "clusterId")
		ClusterID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert code to int", err)
			return
		}
		Cluster, err := cfg.DB.GetCluster(r.Context(), int32(ClusterID))
		if err != nil {
			var msg string
			if err == sql.ErrNoRows {
				msg = "Cluster not found"
			} else {
				msg = "Failed to get cluster"
			}
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), clusterKey, &Cluster)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) ListClusters(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Listing clusters")
	clusters, err := cfg.DB.ListClusters(r.Context())
	if err != nil {
		cfg.Log.Error("Failed to list clusters(ListClusters Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListClusters", err)
		return
	}

	if clusters == nil {
		clusters = []database.Cluster{}
	}

	respondWithJson(w, http.StatusOK, "Clusters retrieved", clusters)
}

func (cfg *Config) GetCluster(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Getting cluster")
	ctx := r.Context()
	cluster, ok := ctx.Value(clusterKey).(*database.Cluster)

	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Cluster not found", nil)
		return
	}

	respondWithJson(w, http.StatusOK, "Cluster retrieved", cluster)
}

func (cfg *Config) CreateCluster(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Creating cluster")

	type parameters struct {
		Name     string `json:"name"`
		ImageUrl string `json:"imageUrl"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request Body", err)
		return
	}
	createdCluster, err := cfg.DB.CreateCluster(r.Context(), database.CreateClusterParams{
		Name:     params.Name,
		ImageUrl: params.ImageUrl,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create cluster", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Cluster created", createdCluster)
}

func (cfg *Config) DeleteCluster(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Deleting cluster")
	ctx := r.Context()
	cluster, ok := ctx.Value(clusterKey).(*database.Cluster)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Cluster not found", nil)
		return
	}

	err := cfg.DB.DeleteCluster(r.Context(), int32(cluster.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete cluster", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Cluster deleted", nil)
}

func (cfg *Config) PatchCluster(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("Patching cluster")
	ctx := r.Context()
	cluster, ok := ctx.Value(clusterKey).(*database.Cluster)
	if !ok {
		respondWithError(w, http.StatusUnprocessableEntity, "Cluster not found", nil)
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
		cluster.Name = params.Name
	}
	if params.ImageUrl != "" {
		cluster.ImageUrl = params.ImageUrl
	}

	err := cfg.DB.PatchCluster(r.Context(), database.PatchClusterParams{
		ID:       cluster.ID,
		ImageUrl: cluster.ImageUrl,
		Name:     cluster.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to patch cluster", err)
		return
	}
	respondWithJson(w, http.StatusOK, "Cluster patched", cluster)
}
