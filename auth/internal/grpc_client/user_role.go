package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpclientmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/response/grpclient"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRoleGrpcClientHandler interface {
	AssignRole(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.ApiResponseUserRole, *response.ErrorResponse)
}

type userRoleClientHandler struct {
	client       pbuserrole.UserRoleServiceClient
	logger       logger.LoggerInterface
	errorHandler middlewares.GRPCErrorHandling
	mapper       grpclientmapper.UserRoleClientResponseMapper
}

func NewUserRoleClienthandler(userRoleClient pbuserrole.UserRoleServiceClient, logger logger.LoggerInterface, errorHandler middlewares.GRPCErrorHandling) UserRoleGrpcClientHandler {
	mapper := grpclientmapper.NewUserRoleClientResponseMapper()

	return &userRoleClientHandler{
		client:       userRoleClient,
		logger:       logger,
		errorHandler: errorHandler,
		mapper:       mapper,
	}
}

func (u *userRoleClientHandler) AssignRole(ctx context.Context, request *requests.CreateUserRoleRequest) (*response.ApiResponseUserRole, *response.ErrorResponse) {
	if errResp := u.errorHandler.ValidateClient(u.client, "role-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuserrole.CreateUserRoleRequest{
		Userid: int32(request.UserId),
		Roleid: int32(request.RoleId),
	}

	res, err := u.client.AssignRole(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}

		u.logger.Error("gRPC client error when assign role", zap.Error(err))

		return nil, u.errorHandler.HandleGRPCError(err, "role-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return u.mapper.ToApiResponseUserRole(res), nil
}
