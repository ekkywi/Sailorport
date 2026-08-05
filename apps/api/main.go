package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ekkywi/sailorport/apps/api/internal/config"
	"github.com/ekkywi/sailorport/apps/api/internal/db"
	"github.com/ekkywi/sailorport/apps/api/internal/handler"
)

func main() {
	cfg := config.Load()

	sqlDB, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx, sqlDB); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Printf("Database OK (SELECT 1)")

	mux := http.NewServeMux()

	health := handler.NewHealthHandler("sailorport-api", cfg.Version)
	echo := handler.NewEchoHandler()

	mux.Handle("/healthz", health)
	mux.Handle("/api/v1/echo", echo)

	addr := ":" + cfg.Port
	log.Printf("Sailorport API (%s) running on http://localhost%s", cfg.AppEnv, addr)
	
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}