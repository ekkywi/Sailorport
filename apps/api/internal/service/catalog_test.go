package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func TestNormalizeCreate_RequiresName(t *testing.T) {
	_, err := normalizeCreate(model.CreateServiceRequest{
		Name:        "   ",
		Description: "x",
		Owner:       "y",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("unexpected message: %v", err)
	}
}

func TestNormalizeCreate_TrimsFields(t *testing.T) {
	req, err := normalizeCreate(model.CreateServiceRequest{
		Name:        "  payments-api  ",
		Description: "  desc  ",
		Owner:       "  team  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "payments-api" || req.Description != "desc" || req.Owner != "team" {
		t.Fatalf("fields not trimmed: %+v", req)
	}
}
