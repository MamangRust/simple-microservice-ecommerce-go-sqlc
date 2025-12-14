package response

type OrderResponse struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type OrderResponseDeleteAt struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	TotalPrice int     `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type ApiResponseOrder struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *OrderResponse `json:"data"`
}

type ApiResponseOrderDeleteAt struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    *OrderResponseDeleteAt `json:"data"`
}

type ApiResponseOrderDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseOrderAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationOrderDeleteAt struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []*OrderResponseDeleteAt `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}

type ApiResponsePaginationOrder struct {
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Data       []*OrderResponse `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}
