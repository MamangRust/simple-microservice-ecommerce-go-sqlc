package productprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	helperproto "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/proto/helpers"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
)

type productCommandProtoMapper struct {
}

func NewProductCommandProtoMapper() ProductCommandProtoMapper {
	return &productCommandProtoMapper{}
}

func (p *productCommandProtoMapper) ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproduct.ApiResponseProduct {
	return &pbproduct.ApiResponseProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProduct(pbResponse),
	}
}

func (p *productCommandProtoMapper) ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproduct.ApiResponseProductDeleteAt {
	return &pbproduct.ApiResponseProductDeleteAt{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProductDeleteAt(pbResponse),
	}
}

func (p *productCommandProtoMapper) ToProtoResponseProductDelete(status string, message string) *pbproduct.ApiResponseProductDelete {
	return &pbproduct.ApiResponseProductDelete{
		Status:  status,
		Message: message,
	}
}

func (p *productCommandProtoMapper) ToProtoResponseProductAll(status string, message string) *pbproduct.ApiResponseProductAll {
	return &pbproduct.ApiResponseProductAll{
		Status:  status,
		Message: message,
	}
}

func (p *productCommandProtoMapper) mapResponseProduct(product *response.ProductResponse) *pbproduct.ProductResponse {
	return &pbproduct.ProductResponse{
		Id:        int32(product.ID),
		Name:      product.Name,
		Price:     int64(product.Price),
		Stock:     int32(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (p *productCommandProtoMapper) mapResponseProductDeleteAt(product *response.ProductResponseDeleteAt) *pbproduct.ProductResponseDeleteAt {
	res := &pbproduct.ProductResponseDeleteAt{
		Id:        int32(product.ID),
		Name:      product.Name,
		Price:     int64(product.Price),
		Stock:     int32(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	if product.DeletedAt != nil {
		res.DeletedAt = helperproto.StringPtrToWrapper(product.DeletedAt)
	}

	return res
}
