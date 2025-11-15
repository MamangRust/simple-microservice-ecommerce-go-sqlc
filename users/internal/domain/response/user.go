package response

type UserResponse struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UserResponseWithPassword struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UserResponseDeleteAt struct {
	DeletedAt  *string `json:"deleted_at"`
	ID         int     `json:"id"`
	FirstName  string  `json:"firstname"`
	LastName   string  `json:"lastname"`
	Email      string  `json:"email"`
	IsVerified bool    `json:"is_verified"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type ApiResponseUser struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *UserResponse `json:"data"`
}

type ApiResponseUserDeleteAt struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    *UserResponseDeleteAt `json:"data"`
}

type ApiResponsesUser struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []*UserResponse `json:"data"`
}

type ApiResponseUserDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseUserAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationUserDeleteAt struct {
	Status     string                  `json:"status"`
	Message    string                  `json:"message"`
	Data       []*UserResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta         `json:"pagination"`
}

type ApiResponsePaginationUser struct {
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Data       []*UserResponse `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}
