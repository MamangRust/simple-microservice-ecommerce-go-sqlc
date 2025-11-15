package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	grpclientmapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/grpcclient"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
	"go.uber.org/zap"
)

type UserRoleGrpcClientHandler interface {
	AssignUserRole(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.ApiResponseUserRole, *response.ErrorResponse)
	UpdateUserRole(ctx context.Context, req *requests.UpdateUserRoleRequest) (*response.ApiResponseUserRole, *response.ErrorResponse)
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
		mapper:       mapper,
		errorHandler: errorHandler,
	}
}

func (u *userRoleClientHandler) AssignUserRole(ctx context.Context, request *requests.CreateUserRoleRequest) (*response.ApiResponseUserRole, *response.ErrorResponse) {
	if errResp := u.errorHandler.ValidateClient(u.client, "role-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuserrole.CreateUserRoleRequest{
		Userid: int32(request.UserId),
		Roleid: int32(request.RoleId),
	}

	res, err := u.client.AssignRole(ctx, req)

	if err != nil {

		u.logger.Error("gRPC client error when assign role", zap.Error(err))

		return nil, u.errorHandler.HandleGRPCError(err, "role-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return u.mapper.ToApiResponseUserRole(res), nil
}

func (u *userRoleClientHandler) UpdateUserRole(ctx context.Context, req *requests.UpdateUserRoleRequest) (*response.ApiResponseUserRole, *response.ErrorResponse) {
	if errResp := u.errorHandler.ValidateClient(u.client, "role-service"); errResp != nil {
		return nil, errResp
	}

	reqPb := &pbuserrole.CreateUserRoleRequest{
		Userid: int32(req.UserId),
		Roleid: int32(req.RoleId),
	}

	res, err := u.client.UpdateRole(ctx, reqPb)

	if err != nil {
		return nil, u.errorHandler.HandleGRPCError(err, "role-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return u.mapper.ToApiResponseUserRole(res), nil
}
