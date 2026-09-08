package model

import "time"

const (
	WorkerStatusOnline   = "online"
	WorkerStatusOffline  = "offline"
	WorkerStatusDraining = "draining"
)

type Worker struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Hostname   string         `json:"hostname"`
	Status     string         `json:"status"`
	Labels     map[string]any `json:"labels"`
	LastSeenAt *time.Time     `json:"last_seen_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type RegisterWorkerRequest struct {
	Name     string         `json:"name"`
	Hostname string         `json:"hostname"`
	Labels   map[string]any `json:"labels"`
}

type HeartbeatRequest struct {
	Status string `json:"status"`
}

type UpdateWorkerLabelsRequest struct {
	Tier         string `json:"tier"`
	Environments string `json:"environments"`
}

type DecommissionWorkerRequest struct{}

type RestoreWorkerRequest struct{}

func IsWorkerStatus(s string) bool {
	switch s {
	case WorkerStatusOnline, WorkerStatusOffline, WorkerStatusDraining:
		return true
	default:
		return false
	}
}

func WorkerAcceptsDeploy(status string) bool {
	return status == WorkerStatusOnline
}
