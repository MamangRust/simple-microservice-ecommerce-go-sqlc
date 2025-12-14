package response

type RoleResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ApiResponseRole struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *RoleResponse `json:"data"`
}
