package model

import "time"

type RuntimeJob struct {
	ID              string    `json:"id"`
	ServiceID       string    `json:"service_id"`
	DeploymentID    string    `json:"deployment_id"`
	ServiceName     string    `json:"service_name"`
	EnvironmentSlug string    `json:"environment_slug"`
	Action          string    `json:"action"`
	Status          string    `json:"status"`
	WorkerID        *string   `json:"worker_id"`
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpdateRuntimeJobRequest struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}
