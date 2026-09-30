package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

func (s *DeploymentsStore) ListFailed(ctx context.Context, ownerUserID string, limit int) ([]model.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ownerUserID = strings.TrimSpace(ownerUserID)

	var (
		rows *sql.Rows
		err error
	)

	if ownerUserID == "" {
		const q = `
			SELECT
				d.id, d.service_id, s.name, e.slug,
				d.error_message, d.created_at, d.updated_at
			FROM deployments d
			JOIN services s ON s.id = d.service_id
			JOIN environments e ON e.id = d.environment_id
			WHERE d.status = 'failed'
			ORDER BY d.updated_at DESC
			LIMIT $1
		`
		rows, err = s.db.QueryContext(ctx, q, limit)
	} else {
		const q = `
			SELECT
				d.id, d.service_id, s.name, e.slug,
				d.error_message, d.created_at, d.updated_at
			FROM deployments d
			JOIN services s ON s.id = d.service_id
			JOIN environments e ON e.id = d.environment_id
			WHERE d.status = 'failed'
				AND s.owner_user_id = $1
			ORDER BY d.updated_at DESC
			LIMIT $2
		`
		rows, err = s.db.QueryContext(ctx, q, ownerUserID, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("ListFailed deployments: %w", err)
	}
	defer rows.Close()

	out := make([]model.Notification, 0)
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(
			&n.ID,
			&n.ServiceID,
			&n.ServiceName,
			&n.EnvironmentSlug,
			&n.ErrorMessage,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		n.Type = "deploy_failed"
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}