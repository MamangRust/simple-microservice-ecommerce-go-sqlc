package requests

import "github.com/go-playground/validator"

type FindAllOrder struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateOrderRecordRequest struct {
	UserID     int `json:"user_id" validate:"required"`
	TotalPrice int `json:"total_price"`
}

type UpdateOrderRecordRequest struct {
	OrderID    int `json:"order_id"`
	UserID     int `json:"user_id" validate:"required"`
	TotalPrice int `json:"total_price" validate:"required"`
}

type CreateOrderRequest struct {
	UserID int                      `json:"user_id" validate:"required"`
	Items  []CreateOrderItemRequest `json:"items" validate:"required"`
}

type UpdateOrderRequest struct {
	OrderID int                      `json:"order_id"`
	UserID  int                      `json:"user_id" validate:"required"`
	Items   []UpdateOrderItemRequest `json:"items" validate:"required"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
	Price     int `json:"price" validate:"required"`
}

type UpdateOrderItemRequest struct {
	OrderItemID int `json:"order_item_id" validate:"required"`
	ProductID   int `json:"product_id" validate:"required"`
	Quantity    int `json:"quantity" validate:"required"`
	Price       int `json:"price" validate:"required"`
}

func (r *CreateOrderRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateOrderRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
