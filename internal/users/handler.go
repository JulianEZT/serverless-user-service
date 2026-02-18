package users

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/JulianEZT/serverless-user-service/internal/httpapi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const usersPathPrefix = "/users/"

// Handler holds dependencies for user HTTP handlers.
type Handler struct {
	svc *Service
}

// NewHandler returns a new Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// CreateUser handles POST /users.
func (h *Handler) CreateUser(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	reqCtx := req.RequestContext
	requestID := reqCtx.RequestID
	requesterSub := extractSub(reqCtx)
	if requesterSub == "" {
		slog.Warn("missing JWT claims", "requestId", requestID)
		return httpapi.ErrorResponse(401, "unauthorized"), nil
	}
	slog.Info("incoming request", "method", "POST", "path", req.RawPath, "requestId", requestID, "requesterSub", requesterSub)

	var in CreateUserInput
	if err := json.Unmarshal([]byte(req.Body), &in); err != nil {
		return httpapi.ErrorResponse(400, "invalid JSON body"), nil
	}
	goCtx := SetRequestID(context.Background(), requestID)
	u, err := h.svc.CreateUser(goCtx, in, requesterSub)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation: ") {
			return httpapi.ErrorResponse(400, strings.TrimPrefix(err.Error(), "validation: ")), nil
		}
		// User already exists (DynamoDB conditional check or mock)
		if errors.Is(err, ErrUserAlreadyExists) || isConditionalCheckErr(err) {
			return httpapi.ErrorResponse(409, "user already exists"), nil
		}
		slog.Error("create user failed", "requestId", requestID, "error", err)
		return httpapi.ErrorResponse(500, "internal server error"), nil
	}
	slog.Info("DynamoDB write result", "requestId", requestID, "userId", u.ID, "action", "Put")
	return httpapi.JSON(201, u), nil
}

// GetUser handles GET /users/{id}. Id is extracted from req.RawPath.
func (h *Handler) GetUser(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	reqCtx := req.RequestContext
	requestID := reqCtx.RequestID
	requesterSub := extractSub(reqCtx)
	if requesterSub == "" {
		slog.Warn("missing JWT claims", "requestId", requestID)
		return httpapi.ErrorResponse(401, "unauthorized"), nil
	}
	slog.Info("incoming request", "method", "GET", "path", req.RawPath, "requestId", requestID, "requesterSub", requesterSub)

	id := extractIDFromPath(req.RawPath)
	if id == "" {
		return httpapi.ErrorResponse(404, "not found"), nil
	}
	goCtx := SetRequestID(context.Background(), requestID)
	u, err := h.svc.GetUser(goCtx, id)
	if err != nil {
		slog.Error("get user failed", "requestId", requestID, "error", err)
		return httpapi.ErrorResponse(500, "internal server error"), nil
	}
	if u == nil {
		return httpapi.ErrorResponse(404, "not found"), nil
	}
	slog.Info("DynamoDB read result", "requestId", requestID, "userId", u.ID, "action", "GetItem")
	return httpapi.JSON(200, u), nil
}

func extractSub(ctx events.APIGatewayV2HTTPRequestContext) string {
	if ctx.Authorizer == nil || ctx.Authorizer.JWT == nil || ctx.Authorizer.JWT.Claims == nil {
		return ""
	}
	return ctx.Authorizer.JWT.Claims["sub"]
}

func extractIDFromPath(rawPath string) string {
	path := strings.TrimSuffix(rawPath, "/")
	if !strings.HasPrefix(path, usersPathPrefix) {
		return ""
	}
	return strings.TrimPrefix(path, usersPathPrefix)
}

func isConditionalCheckErr(err error) bool {
	var ccf *types.ConditionalCheckFailedException
	return errors.As(err, &ccf)
}