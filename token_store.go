package main

import (
	"errors"
	"sync"
)

// TokenEvent represents a single token generation event.
type TokenEvent struct {
	ID    int    // Unique ID for the token event
	Token string // The generated token
}

// TokenEventStore stores token events and allows retrieving them by ID.
type TokenEventStore struct {
	mu     sync.RWMutex // Mutex to protect concurrent access
	events []TokenEvent // Slice to store all token events
	nextID int          // Next ID to assign to a new token event
}

// NewTokenEventStore initializes a new TokenEventStore.
func NewTokenEventStore() *TokenEventStore {
	return &TokenEventStore{
		events: make([]TokenEvent, 0),
		nextID: 0,
	}
}

// AddTokenEvent adds a new token event to the store.
func (store *TokenEventStore) AddTokenEvent(token string) int {
	store.mu.Lock()
	defer store.mu.Unlock()

	event := TokenEvent{
		ID:    store.nextID,
		Token: token,
	}

	store.events = append(store.events, event)
	store.nextID++

	return event.ID
}

// GetTokenEvents returns a slice of token events starting from the given token ID.
func (store *TokenEventStore) GetTokenEvents(fromID int) ([]TokenEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if fromID < 0 || fromID >= store.nextID {
		return nil, errors.New("invalid token ID")
	}

	// Return all events from the fromID to the end
	return store.events[fromID:], nil
}

// GetLatestTokenEvents returns all token events from the store.
func (store *TokenEventStore) GetLatestTokenEvents() []TokenEvent {
	store.mu.RLock()
	defer store.mu.RUnlock()

	return store.events
}
