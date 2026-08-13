package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{db: db}
}

func (s *UsersStore) Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error) {
	const q = `
		INSERT INTO users (email, name, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, role, disabled, created_at, updated_at`

	u, err := scanUser(s.db.QueryRowContext(ctx, q, email, name, passwordHash, role))
	if err != nil {
		if isUniqueViolation(err) {
			return model.User{}, ErrConflict
		}
		return model.User{}, fmt.Errorf("Create user: %w", err)
	}
	return u, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (model.UserRecord, error) {
	const q = `
		SELECT id, email, name, role, disabled, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`

	var rec model.UserRecord
	err := s.db.QueryRowContext(ctx, q, email).Scan(
		&rec.ID,
		&rec.Email,
		&rec.Name,
		&rec.Role,
		&rec.Disabled,
		&rec.PasswordHash,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.UserRecord{}, ErrNotFound
	}
	if err != nil {
		return model.UserRecord{}, fmt.Errorf("Get user by email: %w", err)
	}
	return rec, nil
}

func (s *UsersStore) GetByID(ctx context.Context, id string) (model.User, error) {
	const q = `
		SELECT id, email, name, role, disabled, created_at, updated_at
		FROM users
		WHERE id = $1`

	u, err := scanUser(s.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("Get user by ID: %w", err)
	}
	return u, nil
}

func (s *UsersStore) List(ctx context.Context) ([]model.User, error) {
	const q = `
		SELECT id, email, name, role, disabled, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("List users: %w", err)
	}
	defer rows.Close()

	var out []model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("Scan user: %w", err)
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("List users: %w", err)
	}
	if out == nil {
		out = []model.User{}
	}
	return out, nil
}

func (s *UsersStore) UpdateRole(ctx context.Context, id, role string) (model.User, error) {
	const q = `
		UPDATE users
		SET role = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, name, role, disabled, created_at, updated_at`

	u, err := scanUser(s.db.QueryRowContext(ctx, q, id, role))
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("Update user role: %w", err)
	}
	return u, nil
}

func (s *UsersStore) UpdateDisabled(ctx context.Context, id string, disabled bool) (model.User, error) {
	const q = `
		UPDATE users
		SET disabled = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, name, role, disabled, created_at, updated_at`

	u, err := scanUser(s.db.QueryRowContext(ctx, q, id, disabled))
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("Update user disabled: %w", err)
	}
	return u, nil
}

func (s *UsersStore) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	const q = `
		UPDATE users
		SET password_hash = $2, updated_at = NOW()
		WHERE id = $1
	`

	res, err := s.db.ExecContext(ctx, q, id, passwordHash)
	if err != nil {
		return fmt.Errorf("Update password hash: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanUser(row interface {
	Scan(dest ...any) error
}) (model.User, error) {
	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Role,
		&u.Disabled,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	return u, err
}
