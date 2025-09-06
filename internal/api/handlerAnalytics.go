package rumiapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kaeba0616/es_backend/internal/database"
)

// GET /analytics/combos/clusters
func (cfg *Config) GetClusterCombos(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	start, err := parseRFC3339Ptr(q.Get("start"))
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(q.Get("end"))
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}

	tier := normalizeTier(q.Get("tier"))
	limit := int32(parseIntDefault(q.Get("limit"), 500))
	offset := int32(parseIntDefault(q.Get("offset"), 0))
	minSamples := int32(parseIntDefault(q.Get("minSamples"), 0)) // 0 => SQL 기본 50

	rows, err := cfg.DB.GetTopClusterCombos(r.Context(), database.GetTopClusterCombosParams{
		// $1,$2,$3 = start,end,tier / $4,$5 = limit,offset / $6 = minSamples
		Column1: start,
		Column2: end,
		Column3: tier,
		Limit:   limit,
		Offset:  offset,
		Column6: minSamples,
	})
	if err != nil {
		respondWithError(w, 500, "query failed", err)
		return
	}

	respondWithJson(w, 200, "ok", rows)
}

// GET /analytics/cw/stats
func (cfg *Config) GetCwStats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	start, err := parseRFC3339Ptr(q.Get("start"))
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(q.Get("end"))
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}

	tier := normalizeTier(q.Get("tier"))
	minSamples := int32(parseIntDefault(q.Get("minSamples"), 0)) // 0 => SQL 기본 50

	rows, err := cfg.DB.GetCwStats(r.Context(), database.GetCwStatsParams{
		// $1,$2,$3 = start,end,tier / $4 = minSamples
		Column1: start,
		Column2: end,
		Column3: tier,
		Column4: minSamples,
	})
	if err != nil {
		respondWithError(w, 500, "query failed", err)
		return
	}

	respondWithJson(w, 200, "ok", rows)
}

func (cfg *Config) GetCwStatTop5(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("GetCWStatTop5")

	q := r.URL.Query()

	start, err := parseRFC3339Ptr(q.Get("start"))
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(q.Get("end"))
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}

	tier := normalizeTier(q.Get("tier"))
	minSamples := int32(parseIntDefault(q.Get("minSamples"), 0)) // 0 => SQL 기본 50

	rows, err := cfg.DB.GetCwStatTop5(r.Context(), database.GetCwStatTop5Params{
		// $1,$2,$3,$4,$5 = start,end,tier,minSamples
		Column1: start,
		Column2: end,
		Column3: tier,
		Column4: minSamples,
	})

	if err != nil {
		cfg.Log.Error("Failed to list character weapon stats(GetCwStatTop5 Query)", "error", err)
		respondWithError(w, http.StatusInternalServerError, "DB error GetCwStatTop5", err)
		return
	}
	if rows == nil {
		rows = []database.GetCwStatTop5Row{}
	}

	respondWithJson(w, http.StatusOK, "CW stat top5 retrieved", rows)
}

// 요청 바디
type compBody struct {
	CwIds      []int32 `json:"cw"`
	Start      *string `json:"start"`
	End        *string `json:"end"`
	Tier       *string `json:"tier"`
	MinSamples *int32  `json:"minSamples"`
}

// 응답 타입
type CompAggCW struct {
	CwIds       any     `json:"cw_ids"`
	Samples     int64   `json:"samples"`
	Wins        int64   `json:"wins"`
	WinRate     any     `json:"win_rate"`
	PickRate    any     `json:"pick_rate"`
	AvgMmr      float64 `json:"avg_mmr"`
	AvgSurvival float64 `json:"avg_survival"`
}

type CompAggCluster struct {
	ClusterIds   any     `json:"cluster_ids"`
	ClusterLabel string  `json:"cluster_label"`
	Samples      int64   `json:"samples"`
	Wins         int64   `json:"wins"`
	WinRate      any     `json:"win_rate"`
	PickRate     any     `json:"pick_rate"`
	AvgMmr       float64 `json:"avg_mmr"`
	AvgSurvival  float64 `json:"avg_survival"`
}

type CompMetricsBoth struct {
	ByCW      *CompAggCW      `json:"by_cw,omitempty"`
	ByCluster *CompAggCluster `json:"by_cluster,omitempty"`
}

func (cfg *Config) GetCompMetrics(w http.ResponseWriter, r *http.Request) {
	var p compBody
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, 400, "bad body", err)
		return
	}
	if len(p.CwIds) != 3 {
		respondWithError(w, 400, "exactly 3 cwIds required", nil)
		return
	}

	tier := normalizeTier(*p.Tier)
	start, err := parseRFC3339Ptr(*p.Start)
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(*p.End)
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}
	minSamples := int32(0)
	if p.MinSamples != nil {
		minSamples = *p.MinSamples // 0이면 SQL 기본 50 (정책에 맞게 조절)
	}

	// cluster 기준
	var byCluster *CompAggCluster
	if row, err := cfg.DB.GetCompMetricsBySelectedCWs(
		r.Context(),
		database.GetCompMetricsBySelectedCWsParams{
			Column1: start,      // $1
			Column2: end,        // $2
			Column3: p.CwIds,    // $3
			Column4: tier,       // $4
			Column5: minSamples, // $5
		},
	); err == nil {
		byCluster = &CompAggCluster{
			ClusterIds:   row.ClusterIds,
			ClusterLabel: row.ClusterLabel,
			Samples:      row.Samples,
			Wins:         row.Wins,
			WinRate:      row.WinRate,
			PickRate:     row.PickRate,
			AvgMmr:       row.AvgMmr,
			AvgSurvival:  row.AvgSurvival,
		}
	}

	// 정확히 같은 cw 3개
	var byCW *CompAggCW
	if row, err := cfg.DB.GetCompMetricsByExactCWs(
		r.Context(),
		database.GetCompMetricsByExactCWsParams{
			Column1: start,      // $1
			Column2: end,        // $2
			Column3: p.CwIds,    // $3
			Column4: tier,       // $4
			Column5: minSamples, // $5
		},
	); err == nil {
		byCW = &CompAggCW{
			CwIds:       row.CwIds,
			Samples:     row.Samples,
			Wins:        row.Wins,
			WinRate:     row.WinRate,
			PickRate:    row.PickRate,
			AvgMmr:      row.AvgMmr,
			AvgSurvival: row.AvgSurvival,
		}
	}

	// by_cw 또는 by_cluster 중 하나라도 있으면 OK
	if byCW == nil && byCluster == nil {
		// 둘 다 없으면 200에 빈결과를 주거나 404로 해도 됨. 여기선 빈결과.
		respondWithJson(w, 200, "ok", CompMetricsBoth{})
		return
	}

	respondWithJson(w, 200, "ok", CompMetricsBoth{
		ByCW:      byCW,
		ByCluster: byCluster,
	})
}

func (cfg *Config) GetCwStatsByCw(w http.ResponseWriter, r *http.Request) {
	cwIDstr := chi.URLParam(r, "cwId")
	cwID, err := strconv.Atoi(cwIDstr)
	if err != nil {
		respondWithError(w, 400, "bad cwId", err)
		return
	}
	q := r.URL.Query()
	tier := normalizeTier(q.Get("tier"))

	start, err := parseRFC3339Ptr(q.Get("start")) // 없으면 SQL에서 최근 14일 기본
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(q.Get("end"))
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}
	minSamples := int32(parseIntDefault(q.Get("minSamples"), 0)) // 0 => SQL 기본 50

	rows, err := cfg.DB.GetOneCwStats(r.Context(), database.GetOneCwStatsParams{
		// $1,$2,$3 = start,end,tier / $4 = cwId / $5 = minSamples
		Column1: start,
		Column2: end,
		Column3: tier,
		ID:      int32(cwID),
		Column5: minSamples,
	})
	if err != nil {
		respondWithError(w, 500, "query failed", err)
		return
	}

	respondWithJson(w, 200, "ok", rows)
}

// GET /analytics/cw/{cwId}/trend
func (cfg *Config) GetCwTrend(w http.ResponseWriter, r *http.Request) {
	cwIDstr := chi.URLParam(r, "cwId")
	cwID, err := strconv.Atoi(cwIDstr)
	if err != nil {
		respondWithError(w, 400, "bad cwId", err)
		return
	}

	q := r.URL.Query()
	tier := normalizeTier(q.Get("tier"))

	start, err := parseRFC3339Ptr(q.Get("start")) // 없으면 SQL에서 최근 14일 기본
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(q.Get("end"))
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}
	minSamples := int32(parseIntDefault(q.Get("minSamples"), 0)) // 0 => SQL 기본 50

	rows, err := cfg.DB.GetCwDailyTrend(r.Context(), database.GetCwDailyTrendParams{
		// $1,$2,$3,$4,$5 = start,end,tier,cwId,minSamples
		Column1: start,
		Column2: end,
		Column3: tier,
		Column4: int32(cwID),
		Column5: minSamples,
	})
	if err != nil {
		respondWithError(w, 500, "query failed", err)
		return
	}

	respondWithJson(w, 200, "ok", rows)
}

// 관리용: MV refresh
func (cfg *Config) RefreshMvTrioTeams(w http.ResponseWriter, r *http.Request) {
	if err := cfg.DB.RefreshMvTrioTeams(r.Context()); err != nil {
		respondWithError(w, 500, "query failed", err)
		return
	}
	respondWithJson(w, 200, "ok", nil)
}
