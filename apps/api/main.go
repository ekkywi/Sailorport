package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ekkywi/sailorport/apps/api/internal/catalogapp"
	"github.com/ekkywi/sailorport/apps/api/internal/config"
	"github.com/ekkywi/sailorport/apps/api/internal/db"
	"github.com/ekkywi/sailorport/apps/api/internal/handler"
	"github.com/ekkywi/sailorport/apps/api/internal/migrate"
	"github.com/ekkywi/sailorport/apps/api/internal/secrets"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
	"github.com/ekkywi/sailorport/apps/api/internal/template"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Config invalid: %v", err)
	}

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
	log.Printf("Catalog apps dir: %s", cfg.CatalogAppDir)
	if err := config.EnsureWorkspaceDir(cfg.WorkspaceDir); err != nil {
		log.Fatalf("Workspace setup failed: %v", err)
	}
	log.Printf("Workspace dir writable OK")
	log.Printf("Agent token auth: enabled")

	serviceStore := store.NewServicesStore(sqlDB)
	catalogEnvStore := store.NewCatalogEnvStore(sqlDB)

	secretsKey, err := cfg.SecretsKeyBytes()
	if err != nil {
		log.Fatalf("Config invalid: %v", err)
	}
	var secretsStore secrets.Store
	if len(secretsKey) == 32 {
		secretsStore, err = secrets.NewEncrypted(catalogEnvStore, secretsKey)
		if err != nil {
			log.Fatalf("Secrets store: %v", err)
		}
		log.Printf("Catalog env secrets: encrypted at rest")
	} else {
		secretsStore = secrets.NewPlaintext(catalogEnvStore)
		log.Printf("Catalog env secrets: plaintext (set SAILORPORT_SECRETS_KEY to encrypt)")
	}

	deploymentsStore := store.NewDeploymentsStore(sqlDB)
	catalogAppsReg := catalogapp.NewRegistry(cfg.CatalogAppDir)
	catalog := service.NewCatalog(serviceStore, deploymentsStore, cfg.WorkspaceDir, catalogAppsReg, secretsStore)
	catalogApps := service.NewCatalogApps(catalogAppsReg)
	templates := template.NewRegistry(cfg.TemplatesDir)
	scaffold := service.NewScaffold(catalog, templates, cfg.WorkspaceDir)
	usersStore := store.NewUsersStore(sqlDB)
	authSvc := service.NewAuth(usersStore, cfg.JWTSecret)
	usersSvc := service.NewUsers(usersStore)
	workersStore := store.NewWorkersStore(sqlDB)
	workersSvc := service.NewWorkers(workersStore)
	envsStore := store.NewEnvironmentsStore(sqlDB)
	envsSvc := service.NewEnvironments(envsStore)
	deploymentsSvc := service.NewDeployments(deploymentsStore, envsStore, catalog, workersSvc, secretsStore)
	runtimeStore := store.NewRuntimeStore(sqlDB)
	runtimeSvc := service.NewRuntime(runtimeStore, deploymentsSvc, catalog)
	catalog.SetCleanupEnqueue(runtimeSvc)
	auditStore := store.NewAuditStore(sqlDB)
	auditSvc := service.NewAudit(auditStore)
	catalog.SetAudit(auditSvc)
	usersSvc.SetAudit(auditSvc)
	webhookSvc := service.NewWebhook(catalog, deploymentsSvc)

	router := handler.NewRouter(handler.API{
		Version:      cfg.Version,
		JWTSecret:    cfg.JWTSecret,
		Catalog:      catalog,
		CatalogApps:  catalogApps,
		Scaffold:     scaffold,
		Auth:         authSvc,
		Users:        usersSvc,
		Workers:      workersSvc,
		Deployments:  deploymentsSvc,
		Runtime:      runtimeSvc,
		Environments: envsSvc,
		AgentToken:   cfg.AgentToken,
		Audit:        auditSvc,
		Webhooks:     webhookSvc,
	})

	addr := ":" + cfg.Port
	log.Printf("Sailorport API (%s) running on http://localhost%s", cfg.AppEnv, addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
