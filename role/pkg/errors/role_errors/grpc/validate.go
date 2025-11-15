package rolegrpcerrors

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateRole = response.NewGrpcError("error", "validation failed: invalid create Role request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateRole = response.NewGrpcError("error", "validation failed: invalid update Role request", int(codes.InvalidArgument))
)
