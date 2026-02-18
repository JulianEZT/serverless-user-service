package users

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/JulianEZT/serverless-user-service/pkg/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSPublisher publishes user events to SQS.
type SQSPublisher struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSPublisher returns an SQSPublisher.
func NewSQSPublisher(client *sqs.Client, queueURL string) *SQSPublisher {
	return &SQSPublisher{client: client, queueURL: queueURL}
}

// PublishUserCreated sends a UserCreated event to SQS.
func (p *SQSPublisher) PublishUserCreated(ctx context.Context, payload UserCreatedEventPayload) error {
	now := time.Now().UTC().Format(time.RFC3339)
	ev := events.NewUserCreatedEnvelope(now, events.UserCreatedV1{
		UserID:    payload.UserID,
		Email:     payload.Email,
		Name:      payload.Name,
		CreatedAt: payload.CreatedAt,
		CreatedBy: payload.CreatedBy,
		RequestID: payload.RequestID,
	})
	body, err := events.MarshalEnvelope(ev)
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}
	bodyStr := string(body)
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: &bodyStr,
	})
	if err != nil {
		slog.Error("SQS publish failed", "error", err, "eventType", events.UserCreatedEventType)
		return err
	}
	slog.Info("SQS publish success", "eventType", events.UserCreatedEventType, "userId", payload.UserID)
	return nil
}
