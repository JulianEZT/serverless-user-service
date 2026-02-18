package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/JulianEZT/serverless-user-service/internal/httpapi"
	"github.com/JulianEZT/serverless-user-service/internal/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var router *httpapi.Router

func init() {
	tableName := os.Getenv("USERS_TABLE")
	queueURL := os.Getenv("EVENTS_QUEUE_URL")
	if tableName == "" || queueURL == "" {
		slog.Error("missing required env: USERS_TABLE and EVENTS_QUEUE_URL must be set")
		os.Exit(1)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)
		os.Exit(1)
	}

	ddb := dynamodb.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	repo := users.NewDynamoRepo(ddb, tableName)
	publisher := users.NewSQSPublisher(sqsClient, queueURL)
	svc := users.NewService(repo, publisher)
	h := users.NewHandler(svc)

	router = httpapi.NewRouter()
	router.Register("POST", "/users", h.CreateUser)
	router.Register("GET", "/users/{id}", h.GetUser)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := req.RequestContext.HTTP.Method
	path := req.RawPath
	if path == "" && req.RequestContext.HTTP.Path != "" {
		path = req.RequestContext.HTTP.Path
	}
	if path == "" {
		path = "/"
	}
	h := router.Route(method, path)
	if h == nil {
		return httpapi.ErrorResponse(404, "not found"), nil
	}
	return h(req)
}

func main() {
	lambda.Start(handler)
}
