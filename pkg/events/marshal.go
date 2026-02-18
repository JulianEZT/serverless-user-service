package events

import "encoding/json"

// MarshalEnvelope marshals an envelope to JSON bytes.
func MarshalEnvelope(e Envelope) ([]byte, error) {
	return json.Marshal(e)
}

// NewUserCreatedEnvelope builds an envelope for UserCreatedV1.
func NewUserCreatedEnvelope(occurredAt string, payload UserCreatedV1) Envelope {
	return Envelope{
		EventType:  UserCreatedEventType,
		Version:    UserCreatedV1Version,
		OccurredAt: occurredAt,
		Payload:    payload,
	}
}
