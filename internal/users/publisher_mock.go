package users

import (
	"context"
	"sync"
)

// MockPublisher is an in-memory EventPublisher for tests. It records published payloads.
type MockPublisher struct {
	mu       sync.Mutex
	Published []UserCreatedEventPayload

	// PublishError, if set, makes PublishUserCreated return this error (e.g. to test SQS failure path).
	PublishError error
}

// NewMockPublisher returns a new MockPublisher.
func NewMockPublisher() *MockPublisher {
	return &MockPublisher{Published: nil}
}

// PublishUserCreated appends the payload to Published and returns PublishError if set.
func (m *MockPublisher) PublishUserCreated(ctx context.Context, payload UserCreatedEventPayload) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Published = append(m.Published, payload)
	return m.PublishError
}
