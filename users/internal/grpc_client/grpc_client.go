package grpcclient

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
	"google.golang.org/grpc"
)

type ServiceConnections struct {
	Role *grpc.ClientConn
}

type Deps struct {
	ServiceConnections *ServiceConnections
	Logger             logger.LoggerInterface
}

type HandlerGrpcClient struct {
	RoleClient     RoleGrpcClientHandler
	UserRoleClient UserRoleGrpcClientHandler
}

func NewHandlerGrpcClient(deps *Deps) *HandlerGrpcClient {
	if deps == nil || deps.ServiceConnections == nil {
		return &HandlerGrpcClient{
			RoleClient:     nil,
			UserRoleClient: nil,
		}
	}

	errorHandling := middlewares.NewGRPCErrorHandling(deps.Logger)

	var clientRole pbrole.RoleQueryServiceClient
	var clientUserRole pbuserrole.UserRoleServiceClient
	if deps.ServiceConnections.Role != nil {
		clientRole = pbrole.NewRoleQueryServiceClient(deps.ServiceConnections.Role)
		clientUserRole = pbuserrole.NewUserRoleServiceClient(deps.ServiceConnections.Role)
	}

	clientRolehandler := NewRoleGrpcClientHandler(clientRole, deps.Logger, errorHandling)
	clientUserRoleHandler := NewUserRoleClienthandler(clientUserRole, deps.Logger, errorHandling)

	return &HandlerGrpcClient{
		RoleClient:     clientRolehandler,
		UserRoleClient: clientUserRoleHandler,
	}
}
