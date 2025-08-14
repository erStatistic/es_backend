package data_analysis

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func commandServer(cfg *Config, args ...string) error {
	const port = "8080"

	r := mux.NewRouter()
	// TODO
	// API SERVER
	log.Println("Starting API Server")
	server := &http.Server{
		Handler: r,
		Addr:    ":" + port,
	}

	log.Fatal(server.ListenAndServe())
	return nil
}
