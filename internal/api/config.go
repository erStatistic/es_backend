package rumiapi

import (
	"log/slog"
	"os"
	"time"

	"github.com/kaeba0616/es_backend/internal/database"
	"github.com/kaeba0616/es_backend/internal/erapi"
)

type Config struct {
	DB          *database.Queries
	Log         *slog.Logger
	erapiClient *erapi.Client
}

func NewConfig(db *database.Queries, log *slog.Logger) *Config {
	erapiClient := erapi.NewClient(time.Second*10, os.Getenv("API_KEY"))
	return &Config{
		DB:          db,
		Log:         log,
		erapiClient: &erapiClient,
	}
}
