package response

type UserRoleResponse struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

type ApiResponseUserRole struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *UserRoleResponse `json:"data"`
}
