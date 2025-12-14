package requests

type CreateUserRoleRequest struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}
