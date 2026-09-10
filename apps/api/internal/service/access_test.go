package service

import (
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestCanAccessService_AdminBypass(t *testing.T) {
	err := canAccessService(model.Service{OwnerUserID: "other"}, "me", "admin")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCanAccessService_OwnerOK(t *testing.T) {
	err := canAccessService(model.Service{OwnerUserID: "me"}, "me", "developer")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCanAccessService_OtherForbidden(t *testing.T) {
	err := canAccessService(model.Service{OwnerUserID: "other"}, "me", "developer")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("got %v", err)
	}
}

func TestCanAccessService_EmptyOwnerForbidden(t *testing.T) {
	err := canAccessService(model.Service{OwnerUserID: ""}, "me", "developer")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("got %v", err)
	}
}
