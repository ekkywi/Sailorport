package model

import "time"

// GitHubPushPayload is a minimal subset of a GitHub push webhook body.
type GitHubPushPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		ID            int64  `json:"id"`
		Name          string `json:"name"`
		FullName      string `json:"full_name"`
		CloneURL      string `json:"clone_url"`
		HTMLURL       string `json:"html_url"`
		DefaultBranch string `json:"default_branch"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
}

// WebhookAck is returned to GitHub (and for local curl smoke tests).
type WebhookAck struct {
	Received     bool   `json:"received"`
	Event        string `json:"event"`
	Ignored      bool   `json:"ignored,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Repo         string `json:"repo,omitempty"`
	CloneURL     string `json:"clone_url,omitempty"`
	Branch       string `json:"branch,omitempty"`
	CommitSHA    string `json:"commit_sha,omitempty"`
	Pusher       string `json:"pusher,omitempty"`
	ServiceID    string `json:"service_id,omitempty"`
	DeploymentID string `json:"deployment_id,omitempty"`
	Environment  string `json:"environment,omitempty"`
}

// WebhookDelivery is a recorded GitHub delivery id (dedupe)
type WebhookDelivery struct {
	DeliveryID   string    `json:"delivery_id"`
	ServiceID    string    `json:"service_id,omitempty"`
	DeploymentID string    `json:"deployment_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
