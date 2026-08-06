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
	"github.com/ekkywi/sailorport/apps/api/internal/service"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
	"github.com/ekkywi/sailorport/apps/api/internal/template"
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
	log.Printf("Templates dir: %s", cfg.TemplatesDir)
	log.Printf("Workspace dir: %s", cfg.WorkspaceDir)

	serviceStore := store.NewServicesStore(sqlDB)
	catalog := service.NewCatalog(serviceStore)
	templates := template.NewRegistry(cfg.TemplatesDir)
	scaffold := service.NewScaffold(catalog, templates, cfg.WorkspaceDir)
	usersStore := store.NewUsersStore(sqlDB)
	authSvc := service.NewAuth(usersStore, cfg.JWTSecret)

	router := handler.NewRouter(handler.API{
		Version:   cfg.Version,
		JWTSecret: cfg.JWTSecret,
		Catalog:   catalog,
		Scaffold:  scaffold,
		Auth:      authSvc,
	})

	addr := ":" + cfg.Port
	log.Printf("Sailorport API (%s) running on http://localhost%s", cfg.AppEnv, addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
