package users

// User is the domain model for a user.
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"` // ISO8601
	CreatedBy string `json:"createdBy"` // JWT sub
}

// CreateUserInput is the request body for creating a user.
type CreateUserInput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
