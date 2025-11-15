package response

type ProductResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ProductResponseDeleteAt struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     int     `json:"price"`
	Stock     int     `json:"stock"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type ApiResponseProduct struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *ProductResponse `json:"data"`
}

type ApiResponseProductDeleteAt struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    *ProductResponseDeleteAt `json:"data"`
}

type ApiResponsesProduct struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []*ProductResponse `json:"data"`
}

type ApiResponseProductDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseProductAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationProductDeleteAt struct {
	Status     string                     `json:"status"`
	Message    string                     `json:"message"`
	Data       []*ProductResponseDeleteAt `json:"data"`
	Pagination PaginationMeta             `json:"pagination"`
}

type ApiResponsePaginationProduct struct {
	Status     string             `json:"status"`
	Message    string             `json:"message"`
	Data       []*ProductResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}
