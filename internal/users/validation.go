package users

import (
	"regexp"
	"strings"
)

// Email format: simple check for something@something.tld
var emailRegex = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)

// ValidateCreateInput validates CreateUserInput. Returns a human-readable error message or empty string.
func ValidateCreateInput(in *CreateUserInput) string {
	if in == nil {
		return "request body is required"
	}
	id := strings.TrimSpace(in.ID)
	if id == "" {
		return "id is required"
	}
	email := strings.TrimSpace(in.Email)
	if email == "" {
		return "email is required"
	}
	if !emailRegex.MatchString(email) {
		return "email must be a valid email address"
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return "name is required"
	}
	return ""
}
