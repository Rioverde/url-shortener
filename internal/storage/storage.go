package storage

import "sync"

// MapStorage is an in-memory key-value store of short keys to URLs.
// It uses a map guarded by an RWMutex so it is safe for concurrent use.
type MapStorage struct {
	urls map[string]string
	mu   sync.RWMutex
}

// NewMapStorage returns a new in-memory MapStorage ready to use.
func NewMapStorage() *MapStorage {
	return &MapStorage{
		urls: make(map[string]string),
	}
}

// Put stores the given URL under the given short key, overwriting any existing value.
func (m *MapStorage) Put(key, url string) {
	// Acquire a write lock so concurrent writers do not race on the map.
	m.mu.Lock()
	// Release the lock when the function returns.
	defer m.mu.Unlock()
	// Save the URL under its short key.
	m.urls[key] = url
}

// Get returns the URL stored under the given short key.
// The second return value is false if no URL is associated with the key.
func (m *MapStorage) Get(key string) (string, bool) {
	// Acquire a read lock so multiple readers can run in parallel.
	m.mu.RLock()
	// Release the lock when the function returns.
	defer m.mu.RUnlock()
	// Look up the URL for the given key.
	v, ok := m.urls[key]
	return v, ok
}
