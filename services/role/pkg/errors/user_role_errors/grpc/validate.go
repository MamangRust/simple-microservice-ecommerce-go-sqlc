package userrolegrpcerrors

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateAssignRole     = response.NewGrpcError("error", "validation failed: invalid create Assign Role request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateUserRole = response.NewGrpcError("error", "validation failed: invalid update User Role request", int(codes.InvalidArgument))
)
