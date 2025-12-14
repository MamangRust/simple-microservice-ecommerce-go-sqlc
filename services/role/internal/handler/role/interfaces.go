package rolehandler

import pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"

type RoleQueryHandleGrpc interface {
	pbrole.RoleQueryServiceServer
}

type RoleCommandHandleGrpc interface {
	pbrole.RoleCommandServiceServer
}
