package orderitemgrpcerror

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))
)
