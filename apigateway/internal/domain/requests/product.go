package requests

import "github.com/go-playground/validator"

type FindAllProduct struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required"`
	Stock int    `json:"stock" validate:"required"`
}

type UpdateProductRequest struct {
	ProductID *int   `json:"product_id" `
	Name      string `json:"name" validate:"required"`
	Price     int    `json:"price" validate:"required"`
	Stock     int    `json:"stock" validate:"required"`
}

type UpdateProductStockRequest struct {
	ProductID int `json:"product_id"`
	Stock     int `json:"stock" validate:"required"`
}

func (r *CreateProductRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateProductRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateProductStockRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}
