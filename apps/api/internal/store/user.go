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
		RETURNING id, email, name, role, created_at, updated_at`

	var u model.User
	err := s.db.QueryRowContext(ctx, q, email, name, passwordHash, role).Scan(
		&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
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
		SELECT id, email, name, role, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`

	var rec model.UserRecord
	err := s.db.QueryRowContext(ctx, q, email).Scan(
		&rec.ID, &rec.Email, &rec.Name, &rec.Role, &rec.PasswordHash, &rec.CreatedAt, &rec.UpdatedAt,
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
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	var u model.User
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("Get user by ID: %w", err)
	}
	return u, nil
}
