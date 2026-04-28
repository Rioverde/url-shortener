package sqlite

import (
	sqlite3 "github.com/mattn/go-sqlite3" // initialize the SQLite driver
)

// IsUniqueConstraintError	is a helper function to check if an error is a SQLite unique constraint violation.
func IsUniqueConstraintError(err error) bool {
	sqliteErr, ok := err.(sqlite3.Error)
	return ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
