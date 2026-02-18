package events

// UserCreatedV1 is the versioned payload for a user-created event.
// Used when publishing to SQS for async processing (Lambda B).
type UserCreatedV1 struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	CreatedAt   string `json:"createdAt"`   // ISO8601
	CreatedBy   string `json:"createdBy"`   // JWT sub (requester)
	RequestID   string `json:"requestId,omitempty"`
}
