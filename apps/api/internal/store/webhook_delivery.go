package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type WebhookDeliveriesStore struct {
	db *sql.DB
}

func NewWebhookDeliveriesStore(db *sql.DB) *WebhookDeliveriesStore {
	return &WebhookDeliveriesStore{db: db}
}

func (s *WebhookDeliveriesStore) Exists(ctx context.Context, deliveryID string) (bool, error) {
	deliveryID = strings.TrimSpace(deliveryID)
	if deliveryID == "" {
		return false, fmt.Errorf("delivery_id is required")
	}

	const q = `SELECT 1 FROM webhook_deliveries WHERE delivery_id = $1`
	var n int
	err := s.db.QueryRowContext(ctx, q, deliveryID).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Exists webhook delivery: %w", err)
	}
	return true, nil
}

func (s *WebhookDeliveriesStore) Insert(
	ctx context.Context,
	deliveryID, serviceID, deploymentID string,
) (model.WebhookDelivery, error) {
	deliveryID = strings.TrimSpace(deliveryID)
	if deliveryID == "" {
		return model.WebhookDelivery{}, fmt.Errorf("delivery_id is required")
	}
	var svc any
	if id := strings.TrimSpace(serviceID); id != "" {
		svc = id
	}
	var dep any
	if id := strings.TrimSpace(deploymentID); id != "" {
		dep = id
	}
	const q = `
		INSERT INTO webhook_deliveries (delivery_id, service_id, deployment_id)
		VALUES ($1, $2, $3)
		RETURNING delivery_id, service_id, deployment_id, created_at`
	var (
		out    model.WebhookDelivery
		svcOut sql.NullString
		depOut sql.NullString
	)
	err := s.db.QueryRowContext(ctx, q, deliveryID, svc, dep).Scan(
		&out.DeliveryID,
		&svcOut,
		&depOut,
		&out.CreatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.WebhookDelivery{}, ErrConflict
		}
		return model.WebhookDelivery{}, fmt.Errorf("Insert webhook delivery: %w", err)
	}
	if svcOut.Valid {
		out.ServiceID = svcOut.String
	}
	if depOut.Valid {
		out.DeploymentID = depOut.String
	}
	return out, nil
}
