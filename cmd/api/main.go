package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	rumiapi "github.com/kaeba0616/es_backend/internal/api"
	"github.com/kaeba0616/es_backend/internal/database"

	"github.com/jackc/pgx/v5/pgxpool" // ✅ pgxpool
)

const port = "3333"

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	// ---------- PGX POOL ----------
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("pgx ParseConfig: %v", err)
	}
	// 풀 옵션(원하면 조정)
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.HealthCheckPeriod = 30 * time.Second
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 0

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	// sqlc Queries: pgx/v5로 생성된 코드가 요구하는 DBTX를 pool이 만족합니다.
	dbQueries := database.New(pool)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	apiCfg := rumiapi.NewConfig(dbQueries, logger)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: apiCfg.Routes(),
	}

	log.Println("starting server on port", port)
	fmt.Println("http://localhost:" + port + "/api/v1/characters")

	// ---------- graceful shutdown ----------
	idleConnsClosed := make(chan struct{})
	go func() {
		// Ctrl+C, SIGTERM 처리
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		pool.Close()
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
