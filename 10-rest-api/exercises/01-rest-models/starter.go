package restmodels

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

// Validate checks business constraints on CreateUserRequest and returns a list of field validation errors.
func (r *CreateUserRequest) Validate() []ValidationError {
	// TODO: Implement validation rules for Username, Email, and Age
	return []ValidationError{
		{Field: "TODO", Message: "implement Validate"},
	}
}
