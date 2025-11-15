package requests

type UpdateUserVerifiedRequest struct {
	UserID     int  `json:"user_id"`
	IsVerified bool `json:"is_verified"`
}

type UpdateUserPasswordRequest struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
	VerifiedCode    string `json:"verified_code"`
	IsVerified      bool   `json:"is_verified"`
}
