package requests

import "github.com/go-playground/validator/v10"

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
	VerifiedCode    string `json:"verified_code"`
	IsVerified      bool   `json:"is_verified"`
}

type CreateResetTokenRequest struct {
	UserID     int    `json:"user_id" validate:"required"`
	ResetToken string `json:"reset_token" validate:"required"`
	ExpiredAt  string `json:"expired_at" validate:"required"`
}

type CreateResetPasswordRequest struct {
	ResetToken      string `json:"reset_token" validate:"required"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

func (r *RefreshTokenRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}
func (r *AuthRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r *RegisterRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r *CreateResetPasswordRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r *ForgotPasswordRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

type JWTToken struct {
	Token string `json:"token"`
}
