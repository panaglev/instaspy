package sqlite

import (
	"database/sql"
	"instaspy/src/logger"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type FileInfo struct {
	Id           int
	Username     string
	Hash         string
	Picture_name string
	Havesent     int
}

func New(storagePath string) (*Storage, error) {
	const op = "src.storage.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		logger.HandleOpError(op, err)
		return nil, err
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS info(
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL,
		hash TEXT NOT NULL,
		picture_name TEXT NOT NULL,
		havesent INTEGER NOT NULL);
	`)

	if err != nil {
		logger.HandleOpError(op, err)
		db.Close()
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		logger.HandleOpError(op, err)
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	const op = "src.storage.Close"

	s.Close()
}

func (s *Storage) AddInfo(fileInfo FileInfo) error {
	const op = "src.storage.AddInfo"

	stmt, err := s.db.Prepare("INSERT INTO info(username, hash, picture_name, havesent) VALUES(?, ?, ?, ?)")
	if err != nil {
		logger.HandleOpError(op, err)
		return err
	}

	_, err = stmt.Exec(fileInfo.Username, fileInfo.Hash, fileInfo.Picture_name, 0)
	if err != nil {
		logger.HandleOpError(op, err)
		return err
	}

	return nil
}

func (s *Storage) CheckHash(hash string) (bool, error) {
	const op = "src.storage.CheckHash"

	// Prepare query
	stmtExists, err := s.db.Prepare("SELECT COUNT(*) FROM info WHERE hash = ?")
	if err != nil {
		logger.HandleOpError(op, err)
		return false, err
	}

	// Counter var
	var count int

	// Execute statement
	err = stmtExists.QueryRow(hash).Scan(&count)
	if err != nil {
		return false, err
	}

	// If counter > 0 -> already in db
	if count > 0 {
		return true, nil
		// Else not in db
	} else {
		return false, nil
	}
}
