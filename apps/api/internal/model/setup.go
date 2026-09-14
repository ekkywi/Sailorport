package model

type SetupStatus struct {
	NeedsSetup bool `json:"needs_setup"`
}

type SetupAdminRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
