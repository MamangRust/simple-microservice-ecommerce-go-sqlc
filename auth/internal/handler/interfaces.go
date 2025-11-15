package handler

import (
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
)

type AuthHandleGrpc interface {
	pb.AuthServiceServer
}
