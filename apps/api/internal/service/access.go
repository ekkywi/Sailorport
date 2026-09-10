package service

import (
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func canAccessService(svc model.Service, actorID, role string) error {
	role = strings.TrimSpace(role)
	actorID = strings.TrimSpace(actorID)
	if role == "admin" {
		return nil
	}
	if actorID == "" {
		return fmt.Errorf("%w: missing authenticated user", ErrForbidden)
	}
	ownerID := strings.TrimSpace(svc.OwnerUserID)
	if ownerID == "" {
		return fmt.Errorf("%w: service has no owner", ErrForbidden)
	}
	if ownerID != actorID {
		return fmt.Errorf("%w: not the service owner", ErrForbidden)
	}
	return nil
}