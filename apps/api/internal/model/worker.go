package model

import "time"

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
