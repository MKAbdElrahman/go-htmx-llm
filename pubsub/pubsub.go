package pubsub

import (
	"log"
	"sync"
)

// Subscriber defines the callback function signature for subscribers.
type Subscriber func(payload interface{})

// PubSub is an in-memory publish-subscribe system.
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]Subscriber
}

// NewPubSub creates a new PubSub instance.
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]Subscriber),
	}
}

// Subscribe registers a subscriber to a specific event type.
func (ps *PubSub) Subscribe(eventType string, subscriber Subscriber) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.subscribers[eventType]; !exists {
		ps.subscribers[eventType] = []Subscriber{}
	}
	ps.subscribers[eventType] = append(ps.subscribers[eventType], subscriber)
}

// Publish broadcasts an event to all subscribers of the given event type.
func (ps *PubSub) Publish(eventType string, payload interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subscribers, exists := ps.subscribers[eventType]
	if !exists {
		return
	}

	for _, subscriber := range subscribers {
		go func(sub Subscriber) {
			sub(payload)
		}(subscriber)
	}
}

// Debugging helper to log published events (optional).
func (ps *PubSub) DebugPublish(eventType string, payload interface{}) {
	log.Printf("Publishing event: %s | Payload: %+v\n", eventType, payload)
	ps.Publish(eventType, payload)
}
