package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func (d *Deployments) ListNotifications(ctx context.Context, actorID, role string, limit int) ([]model.Notification, error) {
	role = strings.TrimSpace(role)
	actorID = strings.TrimSpace(actorID)

	switch role {
	case "admin":
		return d.store.ListFailed(ctx, "", limit)
	case "developer":
		if actorID == "" {
			return nil, fmt.Errorf("%w: missing authenticated user", ErrForbidden)
		}
		return d.store.ListFailed(ctx, actorID, limit)
	case "viewer":
		return []model.Notification{}, nil
	default:
		return nil, fmt.Errorf("%w: forbidden", ErrForbidden)
	}
}
