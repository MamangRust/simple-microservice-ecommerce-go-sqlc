package productrepositoryerror

import "errors"

var (
	ErrCreateProduct             = errors.New("failed to create product")
	ErrUpdateProduct             = errors.New("failed to update product")
	ErrUpdateProductCountStock   = errors.New("failed to update product stock count")
	ErrTrashedProduct            = errors.New("failed to move product to trash")
	ErrRestoreProduct            = errors.New("failed to restore product")
	ErrDeleteProductPermanent    = errors.New("failed to permanently delete product")
	ErrRestoreAllProducts        = errors.New("failed to restore all products")
	ErrDeleteAllProductPermanent = errors.New("failed to permanently delete all products")
)
