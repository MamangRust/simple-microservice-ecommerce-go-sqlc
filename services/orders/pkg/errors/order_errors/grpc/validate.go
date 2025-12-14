package ordergrpcerrors

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var ErrGrpcValidateCreateOrder = response.NewGrpcError("error", "validation failed: invalid create order request", int(codes.InvalidArgument))
var ErrGrpcValidateUpdateOrder = response.NewGrpcError("error", "validation failed: invalid update order request", int(codes.InvalidArgument))
