package response

type ProductResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ApiResponseProduct struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *ProductResponse `json:"data"`
}
