package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Rioverde/url-shortener/internal/storage"
	// initialize the SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	// OP is the operation name for error wrapping
	const op = "storage.sqlite.New"
	// Open the SQLite database at the specified path
	db, err := sql.Open("sqlite3", path)
	// If there is an error opening the database, wrap it with the operation name and return the error
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (m *Storage) SaveUrl(key, link string) error {
	const op = "storage.sqlite.SaveUrl"
	// Prepare the SQL statement for inserting a new URL mapping
	stmt, err := m.db.Prepare("INSERT INTO urls (short_code, original_url) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	// Execute the SQL statement with the provided key and URL
	_, err = stmt.Exec(key, link)
	if err != nil {
		// check if the error is a SQLite constraint violation (e.g., duplicate key)
		if IsUniqueConstraintError(err) {
			return fmt.Errorf("%s: key already exists: %w", op, storage.ErrKeyAlreadyExists)
		}
		return fmt.Errorf("%s: failed to execute statement: %w", op, err)
	}

	return nil
}

// GetUrl retrieves the original URL associated with the given key from the database.
func (m *Storage) GetUrl(key string) (string, error) {
	const op = "storage.sqlite.GetUrl"
	// Prepare the SQL statement for retrieving a URL by key
	stmt, err := m.db.Prepare("SELECT original_url FROM urls WHERE short_code = ?")
	if err != nil {
		return "", fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var url string
	// Execute the SQL statement and scan the result into the url variable
	err = stmt.QueryRow(key).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: key not found: %w", op, storage.ErrKeyNotFound)
		}
		return "", fmt.Errorf("%s: failed to execute statement: %w", op, err)
	}

	return url, nil
}

func (m *Storage) DeleteUrl(key string) error {
	const op = "storage.sqlite.DeleteUrl"
	// Prepare the SQL statement for deleting a URL by key
	stmt, err := m.db.Prepare("DELETE FROM urls WHERE short_code = ?")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	// Execute the SQL statement with the provided key
	result, err := stmt.Exec(key)
	if err != nil {
		return fmt.Errorf("%s: failed to execute statement: %w", op, err)
	}

	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	// If there is an error getting the number of affected rows, wrap it with the operation name and return the error
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}
	// If no rows were affected, it means the key was not found, so return a wrapped error indicating that
	if rowsAffected == 0 {
		return fmt.Errorf("%s: key not found: %w", op, storage.ErrNoRowsAffected)
	}

	return nil
}

func (m *Storage) Close() error {
	const op = "storage.sqlite.Close"
	// Close the database connection
	if err := m.db.Close(); err != nil {
		return fmt.Errorf("%s: failed to close database: %w", op, err)
	}
	return nil
}
