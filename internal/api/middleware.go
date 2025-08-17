package rumiapi

import (
	"context"
	"net/http"
	"strconv"
)

func (c *Config) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}

		ctx := context.WithValue(r.Context(), pageKey, page)
		ctx = context.WithValue(ctx, limitKey, limit)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
