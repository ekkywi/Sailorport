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
	if err := config.EnsureWorkspaceDir(cfg.WorkspaceDir); err != nil {
		log.Fatalf("Workspace setup failed: %v", err)
	}
	log.Printf("Workspace dir writable OK")
	log.Printf("Agent token auth: enabled")

	serviceStore := store.NewServicesStore(sqlDB)
	deploymentsStore := store.NewDeploymentsStore(sqlDB)
	catalog := service.NewCatalog(serviceStore, deploymentsStore, cfg.WorkspaceDir)
	templates := template.NewRegistry(cfg.TemplatesDir)
	scaffold := service.NewScaffold(catalog, templates, cfg.WorkspaceDir)
	usersStore := store.NewUsersStore(sqlDB)
	authSvc := service.NewAuth(usersStore, cfg.JWTSecret)
	usersSvc := service.NewUsers(usersStore)
	workersStore := store.NewWorkersStore(sqlDB)
	workersSvc := service.NewWorkers(workersStore)
	envsStore := store.NewEnvironmentsStore(sqlDB)
	envsSvc := service.NewEnvironments(envsStore)
	deploymentsSvc := service.NewDeployments(deploymentsStore, envsStore, catalog, workersSvc)
	runtimeStore := store.NewRuntimeStore(sqlDB)
	runtimeSvc := service.NewRuntime(runtimeStore, deploymentsSvc, catalog)
	catalog.SetCleanupEnqueue(runtimeSvc)
	auditStore := store.NewAuditStore(sqlDB)
	auditSvc := service.NewAudit(auditStore)
	catalog.SetAudit(auditSvc)
	usersSvc.SetAudit(auditSvc)

	router := handler.NewRouter(handler.API{
		Version:      cfg.Version,
		JWTSecret:    cfg.JWTSecret,
		Catalog:      catalog,
		Scaffold:     scaffold,
		Auth:         authSvc,
		Users:        usersSvc,
		Workers:      workersSvc,
		Deployments:  deploymentsSvc,
		Runtime:      runtimeSvc,
		Environments: envsSvc,
		AgentToken:   cfg.AgentToken,
		Audit:        auditSvc,
	})

	addr := ":" + cfg.Port
	log.Printf("Sailorport API (%s) running on http://localhost%s", cfg.AppEnv, addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
