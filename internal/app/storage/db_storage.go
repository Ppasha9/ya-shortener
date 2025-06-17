package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseStorage struct {
	DB *sql.DB
}

func NewDatabase(conn string) (*DatabaseStorage, error) {
	dbConn, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Создаем таблицу shorturls, если ее нет
	_, err = dbConn.Exec(`
		CREATE TABLE IF NOT EXISTS shorturls (
			short_url VARCHAR(8) PRIMARY KEY,
			orig_url TEXT NOT NULL,
		);
	`)

	return &DatabaseStorage{dbConn}, nil
}

func (db *DatabaseStorage) SaveURL(ctx context.Context, shortURL, originalURL string) error {
	_, err := db.DB.ExecContext(ctx, "INSERT INTO shorturls (short_url, orig_url) VALUES (?, ?);", shortURL, originalURL)
	return err
}

func (db *DatabaseStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var origURL string
	err := db.DB.QueryRowContext(ctx, "SELECT orig_url FROM shorturls WHERE short_url = ?;", shortURL).Scan(&origURL)
	if err != nil {
		return "", err
	}

	return origURL, nil
}

func (db *DatabaseStorage) IsExists(ctx context.Context, shortURL string) (bool, error) {
	var exists bool
	err := db.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM shorturls WHERE short_url = ?);", shortURL).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db *DatabaseStorage) Close() error {
	return db.DB.Close()
}
