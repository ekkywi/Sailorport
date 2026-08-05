package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ekkywi/sailorport/apps/api/internal/config"
	"github.com/ekkywi/sailorport/apps/api/internal/db"
	"github.com/ekkywi/sailorport/apps/api/internal/handler"
	"github.com/ekkywi/sailorport/apps/api/internal/migrate"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
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

	if err := migrate.Up(sqlDB); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Printf("Database migrations OK")

	mux := http.NewServeMux()

	health := handler.NewHealthHandler("sailorport-api", cfg.Version)
	echo := handler.NewEchoHandler()
	serviceStore := store.NewServicesStore(sqlDB)
	services := handler.NewServicesHandler(serviceStore)

	mux.Handle("/healthz", health)
	mux.Handle("/api/v1/echo", echo)
	mux.HandleFunc("GET /api/v1/services", services.List)
	mux.HandleFunc("POST /api/v1/services", services.Create)
	mux.HandleFunc("GET /api/v1/services/{id}", services.Get)
	mux.HandleFunc("PUT /api/v1/services/{id}", services.Update)
	mux.HandleFunc("DELETE /api/v1/services/{id}", services.Delete)

	addr := ":" + cfg.Port
	log.Printf("Sailorport API (%s) running on http://localhost%s", cfg.AppEnv, addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}