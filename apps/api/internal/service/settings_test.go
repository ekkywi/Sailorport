package service

import (
	"context"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type fakeSettingsRepo struct {
	open bool
	err  error
}

func (f *fakeSettingsRepo) RegistrationOpen(ctx context.Context) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.open, nil
}

func (f *fakeSettingsRepo) SetRegistrationOpen(ctx context.Context, open bool) error {
	if f.err != nil {
		return f.err
	}
	f.open = open
	return nil
}

func TestSettingsGet(t *testing.T) {
	s := NewSettings(&fakeSettingsRepo{open: true})
	out, err := s.Get(context.Background())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !out.RegistrationOpen {
		t.Fatalf("expected registration_open true")
	}
}

func TestSettingsUpdate_OK(t *testing.T) {
	repo := &fakeSettingsRepo{open: false}
	s := NewSettings(repo)
	open := true
	out, err := s.Update(context.Background(), model.UpdateAppSettingsRequest{
		RegistrationOpen: &open,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !out.RegistrationOpen || !repo.open {
		t.Fatalf("expected open true")
	}
}

func TestSettingsUpdate_MissingField(t *testing.T) {
	s := NewSettings(&fakeSettingsRepo{})
	_, err := s.Update(context.Background(), model.UpdateAppSettingsRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSettingsRegistrationStatus(t *testing.T) {
	s := NewSettings(&fakeSettingsRepo{open: false})
	out, err := s.RegistrationStatus(context.Background())
	if err != nil {
		t.Fatalf("RegistrationStatus: %v", err)
	}
	if out.RegistrationOpen {
		t.Fatalf("expected registration_open false")
	}
}
