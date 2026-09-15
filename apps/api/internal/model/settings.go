package model

import "time"

const SettingKeyRegistrationOpen = "registration_open"

type AppSetting struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppSettings struct {
	RegistrationOpen bool `json:"registration_open"`
}

type UpdateAppSettingsRequest struct {
	RegistrationOpen *bool `json:"registration_open"`
}

type RegistrationStatus struct {
	RegistrationOpen bool `json:"registration_open"`
}


