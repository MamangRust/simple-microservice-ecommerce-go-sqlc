package requests

type UpdateProductStockRequest struct {
	ProductID int `json:"product_id"`
	Stock     int `json:"stock" validate:"required"`
}
