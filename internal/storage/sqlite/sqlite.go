package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"server/internal/storage"
)

type Storage struct {
	db *sql.DB
}

type AliasURL struct {
	url   string
	alias string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    alias TEXT PRIMARY KEY,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx.alias ON url(alias)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s - prepare: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s - exec: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s - exec: %w", op, err)
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s - rows affected: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url.url FROM url WHERE url.alias = ? ")
	if err != nil {
		return "", fmt.Errorf("%s - prepare: %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s - query row: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) (int64, error) {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE url.alias = ?")
	if err != nil {
		return 0, fmt.Errorf("%s - prepare: %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return 0, fmt.Errorf("%s - exec: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s - rows affected: %w", op, err)
	}

	return count, nil
}
