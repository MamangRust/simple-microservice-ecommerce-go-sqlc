package requests

import "github.com/go-playground/validator"

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	ID   *int   `json:"id"`
	Name string `json:"name" validate:"required"`
}

func (r *CreateRoleRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateRoleRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}
