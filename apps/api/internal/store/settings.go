package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type SettingsStore struct {
	db *sql.DB
}

func NewSettingsStore(db *sql.DB) *SettingsStore {
	return &SettingsStore{db: db}
}

func (s *SettingsStore) Get(ctx context.Context, key string) (model.AppSetting, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return model.AppSetting{}, fmt.Errorf("settings key is required")
	}

	const q = `
		SELECT key, value, updated_at
		FROM app_settings
		WHERE key = $1`

	var out model.AppSetting
	err := s.db.QueryRowContext(ctx, q, key).Scan(&out.Key, &out.Value, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AppSetting{}, ErrNotFound
		}
		return model.AppSetting{}, fmt.Errorf("Get setting %q: %w", key, err)
	}
	return out, nil
}

func (s *SettingsStore) Set(ctx context.Context, key, value string) (model.AppSetting, error) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key == "" {
		return model.AppSetting{}, fmt.Errorf("settings key is required")
	}

	const q = `
		INSERT INTO app_settings (key, value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (key) DO UPDATE
		SET value = EXCLUDED.value,
		    updated_at = NOW()
		RETURNING key, value, updated_at`

	var out model.AppSetting
	err := s.db.QueryRowContext(ctx, q, key, value).Scan(&out.Key, &out.Value, &out.UpdatedAt)
	if err != nil {
		return model.AppSetting{}, fmt.Errorf("Set setting %q: %w", key, err)
	}
	return out, nil
}

func (s *SettingsStore) RegistrationOpen(ctx context.Context) (bool, error) {
	st, err := s.Get(ctx, model.SettingKeyRegistrationOpen)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return st.Value == "true", nil
}

func (s *SettingsStore) SetRegistrationOpen(ctx context.Context, open bool) error {
	v := "false"
	if open {
		v = "true"
	}
	_, err := s.Set(ctx, model.SettingKeyRegistrationOpen, v)
	return err
}
