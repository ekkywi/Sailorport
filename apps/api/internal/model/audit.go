package model

import "time"

type AuditEvent struct {
	ID           string         `json:"id"`
	At           time.Time      `json:"at"`
	ActorID      string         `json:"actor_id"`
	ActorEmail   string         `json:"actor_email"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resource_type"`
	ResourceID   string         `json:"resource_id"`
	ResourceName string         `json:"resource_name"`
	Payload      map[string]any `json:"payload"`
}

type AuditRecord struct {
	ActorID      string
	ActorEmail   string
	Action       string
	ResourceType string
	ResourceID   string
	ResourceName string
	Payload      map[string]any
}
