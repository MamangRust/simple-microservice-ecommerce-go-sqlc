package requests

import "github.com/go-playground/validator/v10"

type CreateUserRoleRequest struct {
	UserId int `json:"user_id" validate:"required"`
	RoleId int `json:"role_id" validate:"required"`
}

type UpdateUserRoleRequest struct {
	RoleId int `json:"role_id" validate:"required"`
}

type RemoveUserRoleRequest struct {
	UserId int `json:"user_id" validate:"required"`
	RoleId int `json:"role_id" validate:"required"`
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

func (r *RemoveUserRoleRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}
