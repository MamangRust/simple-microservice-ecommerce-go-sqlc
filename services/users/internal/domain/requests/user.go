package requests

import "github.com/go-playground/validator"

type FindAllUsers struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
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

type CreateUserRequest struct {
	FirstName       string `json:"firstname" validate:"required,alpha"`
	LastName        string `json:"lastname" validate:"required,alpha"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	VerifiedCode    string `json:"verified_code"`
	IsVerified      bool   `json:"is_verified"`
}

type UpdateUserRequest struct {
	UserID          *int   `json:"user_id"`
	FirstName       string `json:"firstname" validate:"required,alpha"`
	LastName        string `json:"lastname" validate:"required,alpha"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UpdateUserVerifiedRequest struct {
	UserID    int  `json:"user_id" validate:"required"`
	IsVerfied bool `json:"is_verified" validate:"required"`
}

type UpdateUserPasswordRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *CreateUserRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateUserRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateUserVerifiedRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

func (r *UpdateUserPasswordRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
