package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kaeba0616/es_backend/internal/erapi"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("API_KEY")

	esClient := erapi.NewClient(time.Second*10, apiKey)
	cfg := &config{
		esapiClient: esClient,
		currentUser: nil,
		users:       []erapi.User{},
	}

	startRepl(cfg)
}
