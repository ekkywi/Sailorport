package model

import "time"

type Notification struct {
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	ServiceID       string    `json:"service_id"`
	ServiceName     string    `json:"service_name"`
	EnvironmentSlug string    `json:"environment_slug"`
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
