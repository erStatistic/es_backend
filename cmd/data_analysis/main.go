package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	data_analysis "github.com/kaeba0616/es_backend/internal/data_analysis"
	"github.com/kaeba0616/es_backend/internal/database"
	"github.com/kaeba0616/es_backend/internal/erapi"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	esClient := erapi.NewClient(time.Second*10, apiKey)
	cfg := data_analysis.Config{
		EsapiClient: esClient,
		CurrentUser: nil,
		Users:       []erapi.User{},
		Rankers:     []erapi.User{},
		Bb:          dbQueries,
	}
	data_analysis.StartRepl(&cfg)
}
