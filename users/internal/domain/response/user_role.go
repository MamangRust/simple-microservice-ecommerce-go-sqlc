package response

type UserRoleResponse struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}

type ApiResponseUserRole struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *UserRoleResponse `json:"data"`
}
