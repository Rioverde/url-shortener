package repo

// error definitions for the storage package
import "errors"

// ErrKeyNotFound is returned by the storage when a requested key does not exist.
// Callers can detect it with errors.Is(err, ErrKeyNotFound).
var ErrKeyNotFound = errors.New("key not found")

// ErrKeyAlreadyExists is returned by the storage when trying to save a key that already exists.
// Callers can detect it with errors.Is(err, ErrKeyAlreadyExists).
var ErrKeyAlreadyExists = errors.New("key already exists")

// ErrNoRowsAffected is returned by the storage when a delete operation does not affect any rows.
// Callers can detect it with errors.Is(err, ErrNoRowsAffected).
var ErrNoRowsAffected = errors.New("no rows affected")
