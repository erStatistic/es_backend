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

// POST /analytics/comp/metrics
type compBody struct {
	CwIds      []int32 `json:"cwIds"`
	Start      string  `json:"start"`
	End        string  `json:"end"`
	Tier       string  `json:"tier"`
	MinSamples *int32  `json:"minSamples,omitempty"` // 없으면 기본(nil) -> 0 처리
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

	tier := normalizeTier(p.Tier)
	start, err := parseRFC3339Ptr(p.Start)
	if err != nil {
		respondWithError(w, 400, "bad start", err)
		return
	}
	end, err := parseRFC3339Ptr(p.End)
	if err != nil {
		respondWithError(w, 400, "bad end", err)
		return
	}

	minSamples := int32(0)
	if p.MinSamples != nil {
		minSamples = *p.MinSamples // 0이면 SQL에서 기본 50로 처리
	}

	res, err := cfg.DB.GetCompMetricsBySelectedCWs(r.Context(), database.GetCompMetricsBySelectedCWsParams{
		// $1,$2 = start,end / $3 = cwIds / $4 = tier / $5 = minSamples
		Column1: start,
		Column2: end,
		Column3: p.CwIds,
		Column4: tier,
		Column5: minSamples,
	})
	if err != nil {
		respondWithError(w, 500, "query failed", err)
		return
	}

	respondWithJson(w, 200, "ok", res)
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
