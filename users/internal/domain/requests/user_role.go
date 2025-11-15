package requests

import "github.com/go-playground/validator"

type CreateUserRoleRequest struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}

type UpdateUserRoleRequest struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}

func (r *CreateUserRoleRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

func (r *UpdateUserRoleRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
