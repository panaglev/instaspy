package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// To store: username, picture hash, picture name, have been sent to telegram
type FileInfo struct {
	Id           int
	Username     string
	Hash         string
	Picture_name int
	Havesent     int
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
		username TEXT NOT NULL,
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

func (s *Storage) AddInfo(fileInfo FileInfo) (int, error) {
	const op = "pkg.storage.AddInfo"

	stmtInsert, err := s.db.Prepare("INSERT INTO info(username, hash, picture_name, havesent) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 400, fmt.Errorf("Error preparing query to insert info at %s: %w", op, err)
	}

	_, err = stmtInsert.Exec(fileInfo.Username, fileInfo.Hash, fileInfo.Picture_name, 0)
	if err != nil {
		return 400, fmt.Errorf("Error executing query to insert info at %s: %w", op, err)
	}

	return 200, nil
}

func (s *Storage) CheckHash(hash string) int {
	const op = "src.storage.ChechHash"

	stmtExists, err := s.db.Prepare("SELECT COUNT(*) FROM info WHERE hash = ?")
	if err != nil {
		return 400
	}

	var count int
	err = stmtExists.QueryRow(hash).Scan(&count)
	if err != nil {
		return 400
	}

	if count > 0 {
		return 409 //Object exists
	} else {
		return 200
	}
}
