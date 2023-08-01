package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "pkg.storage.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS info(
		id SERIAL PRIMARY KEY,
		hash TEXT NOT NULL,
		picture_name INTEGER NOT NULL,
		havesent INTEGER NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddPhoto() error {
	const op = "pkg.storage.AddPhoto"

	return nil
	//stmt, err := s.db.Prepare("INSERT INTO info(hash, picture_name, havesent)")
}

/*
	Add posts to db:
	To store: username, picture hash, picture name, have been sent to instagram
*/
