package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	serviceerrors "github.com/Ppasha9/ya-shortener/internal/app/errors"

	"github.com/Ppasha9/ya-shortener/internal/app/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseStorage struct {
	DB *sql.DB
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.UniqueViolation
	}
	return false
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
	_, err = dbConn.Exec("CREATE TABLE IF NOT EXISTS shorturls (shorturl VARCHAR(8) PRIMARY KEY, origurl TEXT NOT NULL, userid BIGINT NOT NULL);")
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Создаем уникальный индекс на поле с исходным урлом
	_, err = dbConn.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_orig_url_idx ON shorturls (userid, origurl);")
	if err != nil {
		return nil, fmt.Errorf("failed to create unique index: %w", err)
	}

	return &DatabaseStorage{dbConn}, nil
}

func (db *DatabaseStorage) SaveURL(ctx context.Context, userID uint32, shortURL, originalURL string) (string, error) {
	_, err := db.DB.ExecContext(ctx, "INSERT INTO shorturls (shorturl, origurl, userid) VALUES ($1, $2, $3);", shortURL, originalURL, userID)
	if isUniqueViolation(err) {
		shortURL, err = db.GetShortURL(ctx, userID, originalURL)
		if err != nil {
			return shortURL, err
		}
		return shortURL, serviceerrors.ErrOrigURLDuplicate
	}
	return shortURL, err
}

func (db *DatabaseStorage) SaveURLs(ctx context.Context, userID uint32, urls []model.URLsPair) error {
	// начинаем транзакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	for _, v := range urls {
		_, err := tx.ExecContext(ctx, "INSERT INTO shorturls (shorturl, origurl, userid) VALUES ($1, $2, $3);", v.ShortURL, v.OrigURL, userID)
		if err != nil {
			// если ошибка, то откатываем изменения
			tx.Rollback()
			return err
		}
	}

	// завершаем транзакцию
	return tx.Commit()
}

func (db *DatabaseStorage) GetOriginalURL(ctx context.Context, userID uint32, shortURL string) (string, error) {
	var origURL string
	err := db.DB.QueryRowContext(ctx, "SELECT origurl FROM shorturls WHERE shorturl = $1 AND userid = $2;", shortURL, userID).Scan(&origURL)
	if err != nil {
		return "", err
	}

	return origURL, nil
}

func (db *DatabaseStorage) GetUserURLs(ctx context.Context, userID uint32) ([]model.URLsPair, error) {
	res := make([]model.URLsPair, 0)
	rows, err := db.DB.QueryContext(ctx, "SELECT shorturl, origurl FROM shorturls WHERE userid = $1", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return res, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v model.URLsPair
		err = rows.Scan(&v.ShortURL, &v.OrigURL)
		if err != nil {
			return nil, err
		}

		res = append(res, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *DatabaseStorage) GetShortURL(ctx context.Context, userID uint32, origURL string) (string, error) {
	var shortURL string
	err := db.DB.QueryRowContext(ctx, "SELECT shorturl FROM shorturls WHERE origurl = $1 AND userid = $2;", origURL, userID).Scan(&shortURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (db *DatabaseStorage) IsExists(ctx context.Context, userID uint32, shortURL string) (bool, error) {
	var cnt int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM shorturls WHERE shorturl = $1 AND userid = $2;", shortURL, userID).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (db *DatabaseStorage) Close() error {
	return db.DB.Close()
}
