package grpcclient

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
	"google.golang.org/grpc"
)

type ServiceConnections struct {
	User *grpc.ClientConn
	Role *grpc.ClientConn
}

type Deps struct {
	ServiceConnections *ServiceConnections
	Logger             logger.LoggerInterface
}

type HandlerGrpcClient struct {
	UserClient     UserGrpcClientHandler
	RoleClient     RoleGrpcClientHandler
	UserRoleClient UserRoleGrpcClientHandler
}

func NewHandlerGrpcClient(deps *Deps) *HandlerGrpcClient {
	if deps == nil || deps.ServiceConnections == nil {
		return &HandlerGrpcClient{
			UserClient:     nil,
			RoleClient:     nil,
			UserRoleClient: nil,
		}
	}

	errorHandling := middlewares.NewGRPCErrorHandling(deps.Logger)

	var clientUserQuery pbuser.UserQueryServiceClient
	var clientUserCommand pbuser.UserCommandServiceClient
	if deps.ServiceConnections.User != nil {
		clientUserQuery = pbuser.NewUserQueryServiceClient(deps.ServiceConnections.User)
		clientUserCommand = pbuser.NewUserCommandServiceClient(deps.ServiceConnections.User)
	}

	var clientRole pbrole.RoleQueryServiceClient
	var clientUserRole pbuserrole.UserRoleServiceClient
	if deps.ServiceConnections.Role != nil {
		clientRole = pbrole.NewRoleQueryServiceClient(deps.ServiceConnections.Role)
		clientUserRole = pbuserrole.NewUserRoleServiceClient(deps.ServiceConnections.Role)
	}

	clientUserHandler := NewUserGrpcClientHandler(clientUserQuery, clientUserCommand, deps.Logger, errorHandling)
	clientRolehandler := NewRoleGrpcClientHandler(clientRole, deps.Logger, errorHandling)
	clientUserRoleHandler := NewUserRoleClienthandler(clientUserRole, deps.Logger, errorHandling)

	return &HandlerGrpcClient{
		UserClient:     clientUserHandler,
		RoleClient:     clientRolehandler,
		UserRoleClient: clientUserRoleHandler,
	}
}
