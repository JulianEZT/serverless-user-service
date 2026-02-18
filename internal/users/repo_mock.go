package users

import (
	"context"
	"errors"
	"sync"
)

// ErrUserAlreadyExists is returned by MockRepo.Put when the user id already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

// MockRepo is an in-memory UserRepository for tests. It mimics DynamoDB behavior:
// Put fails if the user id already exists (like attribute_not_exists(pk)).
type MockRepo struct {
	mu    sync.RWMutex
	users map[string]*User

	// Optional: inject errors for tests (e.g. simulate DynamoDB/SQS failures)
	PutError    error // if set, Put returns this error
	GetByIDError error // if set, GetByID returns (nil, this error)
}

// NewMockRepo returns a new MockRepo (empty store).
func NewMockRepo() *MockRepo {
	return &MockRepo{users: make(map[string]*User)}
}

// Put stores the user. Returns ErrUserAlreadyExists if id already exists.
func (m *MockRepo) Put(ctx context.Context, u *User) error {
	if m.PutError != nil {
		return m.PutError
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.users[u.ID]; exists {
		return ErrUserAlreadyExists
	}
	// Store a copy so callers can't mutate
	cp := *u
	m.users[u.ID] = &cp
	return nil
}

// GetByID returns the user by id, or nil if not found.
func (m *MockRepo) GetByID(ctx context.Context, id string) (*User, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}
