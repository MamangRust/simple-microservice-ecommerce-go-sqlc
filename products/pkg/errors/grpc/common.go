package producgrpcerrror

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	"google.golang.org/grpc/codes"
)

var ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))
