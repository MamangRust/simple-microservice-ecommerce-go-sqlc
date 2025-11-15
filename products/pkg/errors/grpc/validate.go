package producgrpcerrror

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateProduct = response.NewGrpcError("error", "validation failed: invalid create product request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateProduct = response.NewGrpcError("error", "validation failed: invalid update product request", int(codes.InvalidArgument))
)
