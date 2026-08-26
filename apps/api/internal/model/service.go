package model

import "time"

type Service struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	Description           string                `json:"description"`
	Owner                 string                `json:"owner"`
	TemplateID            string                `json:"template_id"`
	WorkspacePath         string                `json:"workspace_path"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	LatestDeployment      *Deployment           `json:"latest_deployment,omitempty"`
	EnvDeployments        map[string]Deployment `json:"env_deployments,omitempty"`
	SourceType            string                `json:"source_type"`
	RepoURL               string                `json:"repo_url"`
	Branch                string                `json:"branch"`
	DockerfilePath        string                `json:"dockerfile_path"`
	WebhookSecret         string                `json:"webhook_secret"`
	WebhookSecretSet      bool                  `json:"webhook_secret_set"`
	AutoDeployEnabled     bool                  `json:"auto_deploy_enabled"`
	AutoDeployEnvironment string                `json:"auto_deploy_environment"`
	CatalogAppID          string                `json:"catalog_app_id"`
	Image                 string                `json:"image"`
	ContainerPort         int                   `json:"container_port"`
}

type CreateServiceRequest struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Owner                 string `json:"owner"`
	TemplateID            string `json:"template_id"`
	WorkspacePath         string `json:"workspace_path"`
	SourceType            string `json:"source_type"`
	RepoURL               string `json:"repo_url"`
	Branch                string `json:"branch"`
	DockerfilePath        string `json:"dockerfile_path"`
	WebhookSecret         string `json:"webhook_secret"`
	AutoDeployEnabled     bool   `json:"auto_deploy_enabled"`
	AutoDeployEnvironment string `json:"auto_deploy_environment"`
	CatalogAppID          string `json:"catalog_app_id"`
	Image                 string `json:"image"`
	ContainerPort         int    `json:"container_port"`
}

type UpdateServiceRequest struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Owner                 string `json:"owner"`
	SourceType            string `json:"source_type"`
	RepoURL               string `json:"repo_url"`
	Branch                string `json:"branch"`
	DockerfilePath        string `json:"dockerfile_path"`
	WebhookSecret         string `json:"webhook_secret"`
	AutoDeployEnabled     *bool  `json:"auto_deploy_enabled"`
	AutoDeployEnvironment string `json:"auto_deploy_environment"`
	CatalogAppID          string `json:"catalog_app_id"`
	Image                 string `json:"image"`
	ContainerPort         int    `json:"container_port"`
}
