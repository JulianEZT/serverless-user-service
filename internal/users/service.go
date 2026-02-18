package users

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// UserRepository defines persistence for users.
type UserRepository interface {
	Put(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}

// EventPublisher publishes events (e.g. to SQS).
type EventPublisher interface {
	PublishUserCreated(ctx context.Context, payload UserCreatedEventPayload) error
}

// UserCreatedEventPayload is the data needed to publish UserCreated.
type UserCreatedEventPayload struct {
	UserID    string
	Email     string
	Name      string
	CreatedAt string
	CreatedBy string
	RequestID string
}

// Service implements user management use cases.
type Service struct {
	repo     UserRepository
	publisher EventPublisher
}

// NewService returns a new Service.
func NewService(repo UserRepository, publisher EventPublisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

// CreateUser creates a user and publishes an event. Returns the created user or a validation/domain error.
func (s *Service) CreateUser(ctx context.Context, in CreateUserInput, createdBy string) (*User, error) {
	if msg := ValidateCreateInput(&in); msg != "" {
		return nil, fmt.Errorf("validation: %s", msg)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	u := &User{
		ID:        strings.TrimSpace(in.ID),
		Email:     strings.TrimSpace(in.Email),
		Name:      strings.TrimSpace(in.Name),
		CreatedAt: now,
		CreatedBy: createdBy,
	}
	if err := s.repo.Put(ctx, u); err != nil {
		return nil, err
	}
	// Best-effort publish; do not fail the request if SQS fails
	_ = s.publisher.PublishUserCreated(ctx, UserCreatedEventPayload{
		UserID:    u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		RequestID: getRequestID(ctx),
	})
	return u, nil
}

// GetUser returns a user by id or nil if not found.
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// contextKey type for request-scoped values
type contextKey string

const requestIDKey contextKey = "requestId"

// SetRequestID stores requestId in context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func getRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}