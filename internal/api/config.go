package rumiapi

import (
	"log/slog"

	"github.com/kaeba0616/es_backend/internal/database"
)

type Config struct {
	DB  *database.Queries
	Log *slog.Logger
}

func NewConfig(db *database.Queries, log *slog.Logger) *Config {
	return &Config{
		DB:  db,
		Log: log,
	}
}
