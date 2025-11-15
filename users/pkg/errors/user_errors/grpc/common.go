package usergrpcerrors

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var ErrGrpcUserInvalidId = response.NewGrpcError("error", "Invalid User ID", int(codes.NotFound))
