package handler

import pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"

type UserQueryHandleGrpc interface {
	pbuser.UserQueryServiceServer
}

type UserCommandHandleGrpc interface {
	pbuser.UserCommandServiceServer
}
