package model

import "time"

type Service struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Owner            string                `json:"owner"`
	TemplateID       string                `json:"template_id"`
	WorkspacePath    string                `json:"workspace_path"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	LatestDeployment *Deployment           `json:"latest_deployment,omitempty"`
	EnvDeployments   map[string]Deployment `json:"env_deployments,omitempty"`
	SourceType       string                `json:"source_type"`
	RepoURL          string                `json:"repo_url"`
	Branch           string                `json:"branch"`
	DockerfilePath   string                `json:"dockerfile_path"`
}

type CreateServiceRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Owner          string `json:"owner"`
	TemplateID     string `json:"template_id"`
	WorkspacePath  string `json:"workspace_path"`
	SourceType     string `json:"source_type"`
	RepoURL        string `json:"repo_url"`
	Branch         string `json:"branch"`
	DockerfilePath string `json:"dockerfile_path"`
}

type UpdateServiceRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Owner          string `json:"owner"`
	SourceType     string `json:"source_type"`
	RepoURL        string `json:"repo_url"`
	Branch         string `json:"branch"`
	DockerfilePath string `json:"dockerfile_path"`
}
