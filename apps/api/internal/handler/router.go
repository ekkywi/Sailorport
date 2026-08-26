package handler

import (
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type API struct {
	Version      string
	JWTSecret    string
	Catalog      *service.Catalog
	Scaffold     *service.Scaffold
	Auth         *service.Auth
	Users        *service.Users
	Workers      *service.Workers
	Deployments  *service.Deployments
	Runtime      *service.Runtime
	AgentToken   string
	Environments *service.Environments
	Audit        *service.Audit
	Webhooks     *service.Webhook
}

func NewRouter(api API) http.Handler {
	mux := http.NewServeMux()
	secret := api.JWTSecret
	token := api.AgentToken
	currentUser := api.Auth

	health := NewHealthHandler("sailorport-api", api.Version)
	authH := NewAuthHandler(api.Auth)
	services := NewServicesHandler(api.Catalog)
	scaffold := NewScaffoldHandler(api.Scaffold)
	workersH := NewWorkersHandler(api.Workers)
	deploymentsH := NewDeploymentsHandler(api.Deployments)
	runtimeH := NewRuntimeHandler(api.Runtime)
	envsH := NewEnvironmentsHandler(api.Environments)
	webhooksH := NewWebhookHandler(api.Webhooks)

	writer := []string{"developer", "admin"}
	reader := []string{"viewer", "developer", "admin"}
	admin := []string{"admin"}

	mux.Handle("/healthz", health)

	// Public: GitHub cannot send a portal JWT (signature check in Step 20c).
	mux.HandleFunc("POST /api/v1/webhooks/github", webhooksH.GitHub)

	mux.HandleFunc("POST /api/v1/auth/register", authH.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)
	mux.Handle("GET /api/v1/auth/me", withAuth(secret, currentUser, authH.Me))

	usersH := NewUsersHandler(api.Users)
	auditH := NewAuditHandler(api.Audit)
	mux.Handle("GET /api/v1/audit", withRole(secret, currentUser, admin, auditH.List))
	mux.Handle("GET /api/v1/users", withRole(secret, currentUser, admin, usersH.List))
	mux.Handle("POST /api/v1/users", withRole(secret, currentUser, admin, usersH.Create))
	mux.Handle("PATCH /api/v1/users/{id}", withRole(secret, currentUser, admin, usersH.Update))
	mux.Handle("POST /api/v1/users/{id}/password", withRole(secret, currentUser, admin, usersH.ResetPassword))
	mux.Handle("DELETE /api/v1/users/{id}", withRole(secret, currentUser, admin, usersH.Delete))

	mux.Handle("GET /api/v1/services", withRole(secret, currentUser, reader, services.List))
	mux.Handle("POST /api/v1/services", withRole(secret, currentUser, writer, services.Create))
	mux.Handle("GET /api/v1/services/{id}", withRole(secret, currentUser, reader, services.Get))
	mux.Handle("PUT /api/v1/services/{id}", withRole(secret, currentUser, writer, services.Update))
	mux.Handle("DELETE /api/v1/services/{id}", withRole(secret, currentUser, writer, services.Delete))

	mux.Handle("GET /api/v1/templates", withRole(secret, currentUser, reader, scaffold.ListTemplates))
	mux.Handle("POST /api/v1/scaffold", withRole(secret, currentUser, writer, scaffold.Create))

	mux.Handle("POST /api/v1/workers/register", withAgentToken(token, workersH.Register))
	mux.Handle("POST /api/v1/workers/{id}/heartbeat", withAgentToken(token, workersH.Heartbeat))
	mux.Handle("GET /api/v1/workers", withRole(secret, currentUser, reader, workersH.List))

	mux.Handle("GET /api/v1/environments", withRole(secret, currentUser, reader, envsH.List))

	mux.Handle("POST /api/v1/services/{id}/deployments", withRole(secret, currentUser, writer, deploymentsH.Create))
	mux.Handle("GET /api/v1/services/{id}/deployments", withRole(secret, currentUser, reader, deploymentsH.ListByService))
	mux.Handle("GET /api/v1/deployments", withRole(secret, currentUser, reader, deploymentsH.List))
	mux.Handle("GET /api/v1/deployments/{id}", withRole(secret, currentUser, reader, deploymentsH.Get))
	mux.Handle("POST /api/v1/deployments/{id}/redeploy", withRole(secret, currentUser, writer, deploymentsH.Redeploy))
	// Update status deployment hanya lewat route agent di bawah: kalau portal boleh
	// PATCH, riwayat dan git_sha (dasar redeploy) bisa dikarang dari sisi user.

	mux.Handle("POST /api/v1/services/{id}/runtime/stop", withRole(secret, currentUser, writer, runtimeH.Stop))
	mux.Handle("POST /api/v1/services/{id}/runtime/start", withRole(secret, currentUser, writer, runtimeH.Start))
	mux.Handle("POST /api/v1/services/{id}/runtime/logs", withRole(secret, currentUser, reader, runtimeH.Logs))
	mux.Handle("GET /api/v1/runtime/{id}", withRole(secret, currentUser, reader, runtimeH.Get))

	mux.Handle("POST /api/v1/agent/jobs/next", withAgentToken(token, deploymentsH.ClaimNext))
	mux.Handle("PATCH /api/v1/agent/deployments/{id}", withAgentToken(token, deploymentsH.Update))
	mux.Handle("POST /api/v1/agent/runtime/next", withAgentToken(token, runtimeH.ClaimNext))
	mux.Handle("PATCH /api/v1/agent/runtime/{id}", withAgentToken(token, runtimeH.Update))

	return CORS(mux)
}
