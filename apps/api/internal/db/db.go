package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_"github.com/jackc/pgx/v5/stdlib"
)

func Open(databaseURL string) (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Open database: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("Ping database: %w", err)
	}

	return sqlDB, nil
}

func Ping(ctx context.Context, db *sql.DB) error {
	var one int
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("Select 1: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("Unexpected select 1 result: %d", one)
	}
	return nil
}