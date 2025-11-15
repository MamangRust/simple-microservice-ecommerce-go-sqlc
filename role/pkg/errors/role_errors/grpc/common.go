package rolegrpcerrors

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var ErrGrpcRoleInvalidId = response.NewGrpcError("error", "Invalid Role ID", int(codes.NotFound))
