package service

import (
	"context"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type SettingsRepository interface {
	RegistrationOpen(ctx context.Context) (bool, error)
	SetRegistrationOpen(ctx context.Context, open bool) error
}

type Settings struct {
	repo SettingsRepository
}

func NewSettings(repo SettingsRepository) *Settings {
	return &Settings{repo: repo}
}

func (s *Settings) Get(ctx context.Context) (model.AppSettings, error) {
	open, err := s.repo.RegistrationOpen(ctx)
	if err != nil {
		return model.AppSettings{}, fmt.Errorf("get settings: %w", err)
	}
	return model.AppSettings{RegistrationOpen: open}, nil
}

func (s *Settings) Update(ctx context.Context, req model.UpdateAppSettingsRequest) (model.AppSettings, error) {
	if req.RegistrationOpen == nil {
		return model.AppSettings{}, fmt.Errorf(
			"%w: registration_open is required",
			ErrInvalid,
		)
	}
	if err := s.repo.SetRegistrationOpen(ctx, *req.RegistrationOpen); err != nil {
		return model.AppSettings{}, fmt.Errorf("update settings: %w", err)
	}
	return model.AppSettings{RegistrationOpen: *req.RegistrationOpen}, nil
}

func (s *Settings) RegistrationStatus(ctx context.Context) (model.RegistrationStatus, error) {
	open, err := s.repo.RegistrationOpen(ctx)
	if err != nil {
		return model.RegistrationStatus{}, fmt.Errorf("registration status: %w", err)
	}
	return model.RegistrationStatus{RegistrationOpen: open}, nil
}
