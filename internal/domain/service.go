package domain

import (
	"fmt"
)

const sizeOfShortString = 6

// URLRepository defines the storage contract that URLService depends on.
// The interface lives in this (consumer) package so the service is not
// coupled to any specific storage implementation.
type URLRepository interface {
	SaveUrl(key, link string) error
	GetUrl(key string) (string, error)
	DeleteUrl(key string) error
}

// URLService holds the business logic for shortening and resolving URLs.
// It depends on a URLRepository for persistence so the storage backend can be swapped.
type URLService struct {
	repo      URLRepository
	generator CodeGenerator
}

// NewURLService builds a URLService backed by the given repository.
func NewURLService(repo URLRepository, gen CodeGenerator) *URLService {
	return &URLService{
		repo:      repo,
		generator: gen,
	}
}

// Shorten is the main function of the URLService
func (s *URLService) Shorten(longURL string) (string, error) {
	if longURL == "" {
		return "", ErrEmptyURL
	}

	// Genetate the random string with size 6
	key, err := s.generator.GenerateRandomString(sizeOfShortString)
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	// Save to memory the string
	err = s.repo.SaveUrl(key, longURL)
	if err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	return key, nil
}
