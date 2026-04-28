package domain

import "errors"

// ErrEmptyURL is returned by URLService.Shorten when the input URL is empty.
// Callers can detect it with errors.Is(err, ErrEmptyURL).
var ErrEmptyURL = errors.New("url cannot be empty")
