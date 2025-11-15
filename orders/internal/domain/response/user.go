package response

type UserResponse struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ApiResponseUser struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *UserResponse `json:"data"`
}
