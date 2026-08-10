package model

import "time"

type Deployment struct {
	ID           string    `json:"id"`
	ServiceID    string    `json:"service_id"`
	WorkerID     *string   `json:"worker_id"`
	Status       string    `json:"status"`
	ImageTag     string    `json:"image_tag"`
	ContainerID  string    `json:"container_id"`
	Port         *int      `json:"port"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateDeploymentRequest struct{}

type UpdateDeploymentRequest struct {
	Status       string `json:"status"`
	ImageTag     string `json:"image_tag"`
	ContainerID  string `json:"container_id"`
	Port         *int   `json:"port"`
	ErrorMessage string `json:"error_message"`
	WorkerID     string `json:"worker_id"`
}

type DeploymentJob struct {
	Deployment
	ServiceName   string `json:"service_name"`
	WorkspacePath string `json:"workspace_path"`
}
