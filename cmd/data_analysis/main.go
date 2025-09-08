package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	data_analysis "github.com/kaeba0616/es_backend/internal/data_analysis"
	"github.com/kaeba0616/es_backend/internal/erapi"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("API_KEY")

	esClient := erapi.NewClient(time.Second*10, apiKey)
	cfg := data_analysis.Config{
		EsapiClient: esClient,
		CurrentUser: nil,
		Users:       []erapi.User{},
		Rankers:     []erapi.User{},
	}
	data_analysis.StartRepl(&cfg)
}
