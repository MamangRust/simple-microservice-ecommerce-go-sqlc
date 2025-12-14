package handler

import pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"

type ProductQueryHandleGrpc interface {
	pbproduct.ProductQueryServiceServer
}

type ProductCommandHandleGrpc interface {
	pbproduct.ProductCommandServiceServer
}
