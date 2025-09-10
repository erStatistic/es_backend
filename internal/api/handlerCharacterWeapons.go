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

var clusterRoleFallback = map[string]string{
	"A": "딜 브루저", "B": "딜 브루저", "C": "딜 브루저", "D": "딜 브루저", "E": "딜 브루저", "F": "서포터",
	"G": "서포터", "H": "스증 마법사", "I": "스증 마법사", "J": "스증 마법사", "K": "스증 마법사", "L": "암살자",
	"M": "암살자", "N": "탱 브루저", "O": "탱 브루저", "P": "탱커", "Q": "탱커", "R": "평원딜", "S": "평원딜",
	"T": "평원딜", "U": "평원딜",
}

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

		ctx := context.WithValue(r.Context(), characterWeaponKey, &CharacterWeapon)
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

// GET /api/v1/cws/directory
func (cfg *Config) ListCwDirectoryByCluster(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("ListCwDirectoryByCluster")

	rows, err := cfg.DB.ListCwDirectoryByCluster(r.Context())
	if err != nil {
		cfg.Log.Error("ListCwDirectoryByCluster failed", "err", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListCwDirectoryByCluster", err)
		return
	}

	type header struct {
		ClusterID int32  `json:"clusterId"`
		Label     string `json:"label"`
		Role      string `json:"role"`
		Counts    struct {
			Cws        int32 `json:"cws"`
			Characters int32 `json:"characters"`
		} `json:"counts"`
	}

	out := make([]header, 0, len(rows))
	for _, r := range rows {
		role := ""
		if r.Role != nil {
			role = *r.Role
		} else if fb, ok := clusterRoleFallback[r.Label]; ok {
			role = fb
		} else {
			role = "기타"
		}
		h := header{
			ClusterID: r.ClusterID,
			Label:     r.Label,
			Role:      role,
		}
		h.Counts.Cws = int32(r.Cws)
		h.Counts.Characters = int32(r.Characters)
		out = append(out, h)
	}

	w.Header().Set("Cache-Control", "public, max-age=300")
	respondWithJson(w, http.StatusOK, "CW directory by cluster retrieved", out)
}

// GET /api/v1/cws/by-cluster/{clusterId}
func (cfg *Config) ListCwEntriesByCluster(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("ListCwEntriesByCluster")

	idStr := chi.URLParam(r, "clusterId")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || id64 <= 0 {
		respondWithError(w, http.StatusBadRequest, "invalid clusterId", err)
		return
	}
	cid := int32(id64)

	// 클러스터 존재 확인
	cl, err := cfg.DB.GetCluster(r.Context(), cid)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Cluster not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get cluster", err)
		return
	}

	rows, err := cfg.DB.ListCwByClusterID(r.Context(), cid)
	if err != nil {
		cfg.Log.Error("ListCwByClusterID failed", "err", err)
		respondWithError(w, http.StatusInternalServerError, "DB error ListCwByClusterID", err)
		return
	}

	type entry struct {
		CwID      int32 `json:"cwId"`
		Character struct {
			ID       int32  `json:"id"`
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl"` // ✅ string
		} `json:"character"`
		Weapon struct {
			Code     int32  `json:"code"`
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl"` // ✅ string
		} `json:"weapon"`
		Position struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
		} `json:"position"`
	}

	out := struct {
		ClusterID int32   `json:"clusterId"`
		Label     string  `json:"label"`
		Entries   []entry `json:"entries"`
	}{
		ClusterID: cl.ID,
		Label:     cl.Name,
		Entries:   make([]entry, 0, len(rows)),
	}

	for _, r := range rows {
		// 캐릭 이미지는 mini 우선, 비어있으면 full 사용 (둘 다 NOT NULL이라 빈문자일 수만 체크)
		charImg := r.ChImgMini
		if charImg == "" && r.ChImgFull != "" {
			charImg = r.ChImgFull
		}

		e := entry{
			CwID: r.CwID,
			Character: struct {
				ID       int32  `json:"id"`
				Name     string `json:"name"`
				ImageURL string `json:"imageUrl"`
			}{
				ID:       r.ChID,
				Name:     r.ChName,
				ImageURL: charImg,
			},
			Weapon: struct {
				Code     int32  `json:"code"`
				Name     string `json:"name"`
				ImageURL string `json:"imageUrl"`
			}{
				Code:     r.WCode,
				Name:     r.WName,
				ImageURL: r.WImg, // 바로 대입
			},
			Position: struct {
				ID   int32  `json:"id"`
				Name string `json:"name"`
			}{
				ID:   r.PID,
				Name: r.PName,
			},
		}

		out.Entries = append(out.Entries, e)
	}

	w.Header().Set("Cache-Control", "public, max-age=300")
	respondWithJson(w, http.StatusOK, "CW entries by cluster retrieved", out)
}

func (cfg *Config) ListCwEntriesByClusters(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("ListCwEntriesByClusters")

	ids, err := parseInt32List(r.URL.Query().Get("ids"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ids", err)
		return
	}
	if len(ids) == 0 {
		respondWithJson(w, http.StatusOK, "empty", map[string]any{"data": []any{}})
		return
	}

	rows, err := cfg.DB.ListCwByClusters(r.Context(), ids)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error ListCwByClusters", err)
		return
	}

	// sqlc가 만든 Row 타입 예: db.ListCwByClustersRow
	type entry struct {
		CwID      int32 `json:"cwId"`
		Character struct {
			ID       int32  `json:"id"`
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl"`
		} `json:"character"`
		Weapon struct {
			Code     *int32  `json:"code,omitempty"`
			Name     *string `json:"name,omitempty"`
			ImageURL *string `json:"imageUrl,omitempty"`
		} `json:"weapon"`
	}
	type bucket struct {
		ClusterID int32   `json:"clusterId"`
		Label     string  `json:"label"`
		Entries   []entry `json:"entries"`
	}

	byCID := make(map[int32]*bucket)
	order := make([]int32, 0, len(ids))
	seenOrder := make(map[int32]bool)

	for _, r := range rows {
		b := byCID[r.ClusterID]
		if b == nil {
			b = &bucket{
				ClusterID: r.ClusterID,
				Label:     r.ClusterLabel,
				Entries:   make([]entry, 0, 32),
			}
			byCID[r.ClusterID] = b
			if !seenOrder[r.ClusterID] {
				order = append(order, r.ClusterID)
				seenOrder[r.ClusterID] = true
			}
		}

		b.Entries = append(b.Entries, entry{
			CwID: r.CwID,
			Character: struct {
				ID       int32  `json:"id"`
				Name     string `json:"name"`
				ImageURL string `json:"imageUrl"`
			}{
				ID:       r.ChID,
				Name:     r.ChName,
				ImageURL: r.ChImg, // COALESCE로 non-null
			},
			Weapon: struct {
				Code     *int32  `json:"code,omitempty"`
				Name     *string `json:"name,omitempty"`
				ImageURL *string `json:"imageUrl,omitempty"`
			}{
				Code:     &r.WCode,
				Name:     &r.WName,
				ImageURL: &r.WImg,
			},
		})
	}

	// 결과를 요청한 ids 순서로 정렬
	out := make([]bucket, 0, len(byCID))
	for _, id := range order {
		if b, ok := byCID[id]; ok {
			out = append(out, *b)
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=300")
	respondWithJson(w, http.StatusOK, "CW entries by clusters retrieved", map[string]any{
		"data": out,
	})
}

type cwOverviewResp struct {
	CwID      int32  `json:"cwId"`
	Tier      string `json:"tier"`
	Character *struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl,omitempty"`
	} `json:"character,omitempty"`
	Weapon *struct {
		Code     int32  `json:"code"`
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl,omitempty"`
	} `json:"weapon,omitempty"`
	Position *struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"position,omitempty"`

	Cluster *struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"cluster,omitempty"`

	Overview struct {
		// summary는 아직 집계테이블 없으니 일단 스텁(0 / 빈배열)
		Summary struct {
			Games       int     `json:"games"`
			WinRate     any     `json:"winRate"`
			PickRate    any     `json:"pickRate"`
			MMRGain     any     `json:"mmrGain"`
			SurvivalSec float64 `json:"survivalSec"`
		} `json:"summary"`
		// ✅ 여기 stats를 character_weapon_stats 에서 채운다
		Stats struct {
			ATK int `json:"atk"`
			DEF int `json:"def"`
			CC  int `json:"cc"`
			SPD int `json:"spd"`
			SUP int `json:"sup"`
		} `json:"stats"`
		Routes []struct {
			ID    int32  `json:"id"`
			Title string `json:"title"`
		} `json:"routes"`
	} `json:"overview"`
}

func (cfg *Config) GetCwOverview(w http.ResponseWriter, r *http.Request) {
	cfg.Log.Info("GetCWOverview")

	idStr := chi.URLParam(r, "cwId")
	cwID64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || cwID64 <= 0 {
		respondWithError(w, http.StatusBadRequest, "invalid cwId", err)
		return
	}
	cwID := int32(cwID64)

	ctx := r.Context()

	// 1) CW 식별(캐릭/무기/포지션)
	ident, err := cfg.DB.GetCwIdentity(ctx, cwID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "CW not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "DB error GetCwIdent", err)
		return
	}

	// 2) 스탯 조회 (없으면 0으로 기본값)
	stats, err := cfg.DB.GetCharacterWeaponStat(ctx, cwID)
	var atk, defv, cc, spd, sup int32
	if err == nil {
		atk, defv, cc, spd, sup = stats.Atk, stats.Def, stats.Cc, stats.Spd, stats.Sup
	} else if err != sql.ErrNoRows {
		respondWithError(w, http.StatusInternalServerError, "DB error GetCwStats", err)
		return
	}
	tier, err := cfg.DB.GetCwTierByCwID(r.Context(), database.GetCwTierByCwIDParams{
		ID: cwID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error GetCwTierByCwID", err)
		return
	}

	// 3) 응답 조립
	out := cwOverviewResp{
		CwID: cwID,
		Tier: tier,
	}

	out.Character = &struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl,omitempty"`
	}{
		ID:       ident.ChID,
		Name:     ident.ChName,
		ImageURL: ident.ChImg,
	}

	out.Weapon = &struct {
		Code     int32  `json:"code"`
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl,omitempty"`
	}{
		Code:     ident.WCode,
		Name:     ident.WName,
		ImageURL: ident.WImg,
	}

	out.Position = &struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   ident.PID,
		Name: ident.PName,
	}

	out.Cluster = &struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   ident.ClusterID,
		Name: ident.ClusterName,
	}

	summary, err := cfg.DB.GetOneCwStats(r.Context(), database.GetOneCwStatsParams{
		ID: cwID,
	})
	fmt.Println(summary)

	if err != nil {
		out.Overview.Summary.Games = 0
		out.Overview.Summary.WinRate = 0
		out.Overview.Summary.PickRate = 0
		out.Overview.Summary.MMRGain = 0
		out.Overview.Summary.SurvivalSec = 0
	} else {
		out.Overview.Summary.Games = int(summary.Samples)
		out.Overview.Summary.WinRate = summary.WinRate
		out.Overview.Summary.PickRate = summary.PickRate
		out.Overview.Summary.MMRGain = summary.AvgMmr
		out.Overview.Summary.SurvivalSec = summary.AvgSurvival
	}
	out.Overview.Stats.ATK = int(atk)
	out.Overview.Stats.DEF = int(defv)
	out.Overview.Stats.CC = int(cc)
	out.Overview.Stats.SPD = int(spd)
	out.Overview.Stats.SUP = int(sup)
	// ✅ routes는 DB 값으로 채움
	routes, err := cfg.DB.ListCWRoutes(r.Context(), database.ListCWRoutesParams{
		CharacterID: ident.ChID,
		WeaponID:    ident.WCode,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list CW routes", err)
		return
	}
	out.Overview.Routes = make([]struct {
		ID    int32  `json:"id"`
		Title string `json:"title"`
	}, len(routes))
	for i, route := range routes {
		out.Overview.Routes[i].ID = route.RouteID
		out.Overview.Routes[i].Title = route.Title
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
	respondWithJson(w, http.StatusOK, "CW overview retrieved", out)
}
