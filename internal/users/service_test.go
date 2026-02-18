package users

import (
	"context"
	"errors"
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	ctx := SetRequestID(context.Background(), "req-1")
	repo := NewMockRepo()
	pub := NewMockPublisher()
	svc := NewService(repo, pub)

	in := CreateUserInput{ID: "u1", Email: "a@b.com", Name: "Alice"}
	u, err := svc.CreateUser(ctx, in, "sub-123")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if u.ID != "u1" || u.Email != "a@b.com" || u.Name != "Alice" || u.CreatedBy != "sub-123" {
		t.Errorf("unexpected user: %+v", u)
	}
	if u.CreatedAt == "" {
		t.Error("CreatedAt should be set")
	}
	if len(pub.Published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(pub.Published))
	}
	if pub.Published[0].UserID != "u1" || pub.Published[0].CreatedBy != "sub-123" || pub.Published[0].RequestID != "req-1" {
		t.Errorf("unexpected payload: %+v", pub.Published[0])
	}

	// Get back from repo
	got, err := svc.GetUser(ctx, "u1")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got == nil || got.ID != "u1" || got.Email != "a@b.com" {
		t.Errorf("GetUser: got %+v", got)
	}
}

func TestService_CreateUser_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepo()
	svc := NewService(repo, NewMockPublisher())

	in := CreateUserInput{ID: "u1", Email: "a@b.com", Name: "Alice"}
	_, err := svc.CreateUser(ctx, in, "sub-1")
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err = svc.CreateUser(ctx, in, "sub-1")
	if err == nil {
		t.Fatal("expected error on second create")
	}
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestService_CreateUser_Validation(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMockRepo(), NewMockPublisher())

	_, err := svc.CreateUser(ctx, CreateUserInput{ID: "", Email: "a@b.com", Name: "A"}, "sub")
	if err == nil || err.Error() != "validation: id is required" {
		t.Errorf("expected validation error, got %v", err)
	}
	_, err = svc.CreateUser(ctx, CreateUserInput{ID: "u1", Email: "invalid", Name: "A"}, "sub")
	if err == nil || err.Error() != "validation: email must be a valid email address" {
		t.Errorf("expected email validation error, got %v", err)
	}
	_, err = svc.CreateUser(ctx, CreateUserInput{ID: "u1", Email: "a@b.com", Name: ""}, "sub")
	if err == nil || err.Error() != "validation: name is required" {
		t.Errorf("expected name validation error, got %v", err)
	}
}

func TestService_GetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMockRepo(), NewMockPublisher())

	u, err := svc.GetUser(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if u != nil {
		t.Errorf("expected nil, got %+v", u)
	}
}

func TestService_CreateUser_PublishErrorDoesNotFailRequest(t *testing.T) {
	ctx := SetRequestID(context.Background(), "req-2")
	repo := NewMockRepo()
	pub := NewMockPublisher()
	pub.PublishError = errors.New("SQS unavailable")
	svc := NewService(repo, pub)

	in := CreateUserInput{ID: "u2", Email: "b@c.com", Name: "Bob"}
	u, err := svc.CreateUser(ctx, in, "sub-2")
	if err != nil {
		t.Fatalf("CreateUser should succeed despite publish error: %v", err)
	}
	if u.ID != "u2" {
		t.Errorf("unexpected user: %+v", u)
	}
	// User should still be in repo
	got, _ := svc.GetUser(ctx, "u2")
	if got == nil {
		t.Error("user should be persisted even when SQS fails")
	}
}
