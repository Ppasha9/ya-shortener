package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func OpenDB(conn string) (*Database, error) {
	dbConn, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	return &Database{dbConn}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
