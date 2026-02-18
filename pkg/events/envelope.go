package events

// Envelope wraps an event with metadata for routing and versioning.
type Envelope struct {
	EventType  string      `json:"eventType"`
	Version    string      `json:"version"`
	OccurredAt string      `json:"occurredAt"` // ISO8601
	Payload    interface{} `json:"payload"`
}

// UserCreatedEventType is the event type string for user-created events.
const UserCreatedEventType = "user.created"

// UserCreatedV1Version is the schema version for UserCreatedV1.
const UserCreatedV1Version = "1"
