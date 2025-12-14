package ordergrpcerrors

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var ErrGrpcFailedInvalidId = response.NewGrpcError("error", "Invalid ID", int(codes.InvalidArgument))
