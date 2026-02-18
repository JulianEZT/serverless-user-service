package httpapi

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// Handler is the signature for a route handler.
type Handler func(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

// Router dispatches by method and path (RawPath). Path params are not parsed here;
// the handler receives the full request and can extract {id} from the path.
type Router struct {
	routes map[string]Handler // key: "METHOD /path" e.g. "POST /users"
}

// NewRouter returns a new Router.
func NewRouter() *Router {
	return &Router{routes: make(map[string]Handler)}
}

// Register associates a handler with method and path.
// Path should be the raw path pattern, e.g. "/users" or "/users/{id}".
// For path params, the handler is responsible for extracting the id from RawPath.
func (r *Router) Register(method, path string, h Handler) {
	r.routes[method+" "+path] = h
}

// Route returns the handler for the given method and rawPath, or nil if not found.
// Path matching: exact match for "/users"; for "/users/{id}" we match prefix "/users/" (handler extracts id from path).
func (r *Router) Route(method, rawPath string) Handler {
	path := strings.TrimSuffix(rawPath, "/")
	if path == "" {
		path = "/"
	}

	key := method + " " + path
	if h, ok := r.routes[key]; ok {
		return h
	}
	// GET /users/{id}
	const usersPrefix = "/users/"
	if strings.HasPrefix(path, usersPrefix) && method == "GET" {
		if h, ok := r.routes["GET /users/{id}"]; ok {
			return h
		}
	}
	return nil
}
