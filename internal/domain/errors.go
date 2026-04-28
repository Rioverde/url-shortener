package domain

import "errors"

// ErrEmptyURL is returned by URLService.Shorten when the input URL is empty.
// Callers can detect it with errors.Is(err, ErrEmptyURL).
var ErrEmptyURL = errors.New("url cannot be empty")

// ErrEmptyKey is returned by URLService.GetUrl when the input key is empty.
// Callers can detect it with errors.Is(err, ErrEmptyKey).
var ErrEmptyKey = errors.New("key cannot be empty")
