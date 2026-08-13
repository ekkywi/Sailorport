package handler

import (
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type API struct {
	Version     string
	JWTSecret   string
	Catalog     *service.Catalog
	Scaffold    *service.Scaffold
	Auth        *service.Auth
	Users       *service.Users
	Workers     *service.Workers
	Deployments *service.Deployments
	Runtime     *service.Runtime
	AgentToken  string
}

func NewRouter(api API) http.Handler {
	mux := http.NewServeMux()
	secret := api.JWTSecret
	token := api.AgentToken

	health := NewHealthHandler("sailorport-api", api.Version)
	authH := NewAuthHandler(api.Auth)
	services := NewServicesHandler(api.Catalog)
	scaffold := NewScaffoldHandler(api.Scaffold)
	workersH := NewWorkersHandler(api.Workers)
	deploymentsH := NewDeploymentsHandler(api.Deployments)
	runtimeH := NewRuntimeHandler(api.Runtime)

	writer := []string{"developer", "admin"}
	reader := []string{"viewer", "developer", "admin"}
	admin := []string{"admin"}

	mux.Handle("/healthz", health)

	mux.HandleFunc("POST /api/v1/auth/register", authH.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)
	mux.Handle("GET /api/v1/auth/me", withAuth(secret, authH.Me))

	usersH := NewUsersHandler(api.Users)
	mux.Handle("GET /api/v1/users", withRole(secret, admin, usersH.List))
	mux.Handle("POST /api/v1/users", withRole(secret, admin, usersH.Create))
	mux.Handle("PATCH /api/v1/users/{id}", withRole(secret, admin, usersH.Update))

	mux.Handle("GET /api/v1/services", withRole(secret, reader, services.List))
	mux.Handle("POST /api/v1/services", withRole(secret, writer, services.Create))
	mux.Handle("GET /api/v1/services/{id}", withRole(secret, reader, services.Get))
	mux.Handle("PUT /api/v1/services/{id}", withRole(secret, writer, services.Update))
	mux.Handle("DELETE /api/v1/services/{id}", withRole(secret, writer, services.Delete))

	mux.Handle("GET /api/v1/templates", withRole(secret, reader, scaffold.ListTemplates))
	mux.Handle("POST /api/v1/scaffold", withRole(secret, writer, scaffold.Create))

	mux.Handle("POST /api/v1/workers/register", withAgentToken(token, workersH.Register))
	mux.Handle("POST /api/v1/workers/{id}/heartbeat", withAgentToken(token, workersH.Heartbeat))
	mux.Handle("GET /api/v1/workers", withRole(secret, reader, workersH.List))

	mux.Handle("POST /api/v1/services/{id}/deployments", withRole(secret, writer, deploymentsH.Create))
	mux.Handle("GET /api/v1/services/{id}/deployments", withRole(secret, reader, deploymentsH.ListByService))
	mux.Handle("GET /api/v1/deployments", withRole(secret, reader, deploymentsH.List))
	mux.Handle("GET /api/v1/deployments/{id}", withRole(secret, reader, deploymentsH.Get))
	mux.Handle("PATCH /api/v1/deployments/{id}", withRole(secret, writer, deploymentsH.Update))

	mux.Handle("POST /api/v1/services/{id}/runtime/stop", withRole(secret, writer, runtimeH.Stop))
	mux.Handle("POST /api/v1/services/{id}/runtime/start", withRole(secret, writer, runtimeH.Start))

	mux.Handle("POST /api/v1/agent/jobs/next", withAgentToken(token, deploymentsH.ClaimNext))
	mux.Handle("PATCH /api/v1/agent/deployments/{id}", withAgentToken(token, deploymentsH.Update))
	mux.Handle("POST /api/v1/agent/runtime/next", withAgentToken(token, runtimeH.ClaimNext))
	mux.Handle("PATCH /api/v1/agent/runtime/{id}", withAgentToken(token, runtimeH.Update))

	return CORS(mux)
}
