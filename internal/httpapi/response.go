package httpapi

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

const contentTypeJSON = "application/json"

// JSON writes statusCode and a JSON body with standard headers.
func JSON(statusCode int, body interface{}) events.APIGatewayV2HTTPResponse {
	raw, err := json.Marshal(body)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		return JSON(statusCode, map[string]string{"error": "internal error"})
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": contentTypeJSON,
		},
		Body: string(raw),
	}
}

// ErrorResponse returns a JSON error response with the given status and message.
func ErrorResponse(statusCode int, message string) events.APIGatewayV2HTTPResponse {
	return JSON(statusCode, map[string]string{"error": message})
}
