package productrepositoryerror

import "errors"

var (
	ErrFindAllProducts = errors.New("failed to find all products")
	ErrFindByActive    = errors.New("failed to find active products")
	ErrFindByTrashed   = errors.New("failed to find trashed products")
	ErrFindByMerchant  = errors.New("failed to find products by merchant")
	ErrFindByCategory  = errors.New("failed to find products by category")
	ErrFindById        = errors.New("failed to find product by ID")
	ErrFindByIdTrashed = errors.New("failed to find trashed product by ID")
)
