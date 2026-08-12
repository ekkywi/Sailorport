package model

import "time"

type Service struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Owner            string      `json:"owner"`
	TemplateID       string      `json:"template_id"`
	WorkspacePath    string      `json:"workspace_path"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	LatestDeployment *Deployment `json:"latest_deployment,omitempty"`
}

type CreateServiceRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Owner         string `json:"owner"`
	TemplateID    string `json:"template_id"`
	WorkspacePath string `json:"workspace_path"`
}

type UpdateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
}
