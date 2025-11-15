package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	grpclientmapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/grpcclient"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	"go.uber.org/zap"
)

type RoleGrpcClientHandler interface {
	FindByName(ctx context.Context, name string) (*response.ApiResponseRole, *response.ErrorResponse)
}

type roleGrpcClientHandler struct {
	client       pbrole.RoleQueryServiceClient
	logger       logger.LoggerInterface
	errorHandler middlewares.GRPCErrorHandling
	mapper       grpclientmapper.RoleClientResponseMapper
}

func NewRoleGrpcClientHandler(roleClient pbrole.RoleQueryServiceClient, logger logger.LoggerInterface, errorHandler middlewares.GRPCErrorHandling) RoleGrpcClientHandler {
	mapper := grpclientmapper.NewRoleClientResponseMapper()

	return &roleGrpcClientHandler{
		client:       roleClient,
		logger:       logger,
		errorHandler: errorHandler,
		mapper:       mapper,
	}
}

func (r *roleGrpcClientHandler) FindByName(ctx context.Context, name string) (*response.ApiResponseRole, *response.ErrorResponse) {
	if errResp := r.errorHandler.ValidateClient(r.client, "role-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbrole.FindByNameRequest{
		Name: name,
	}

	res, err := r.client.FindByName(ctx, req)

	if err != nil {
		return nil, r.errorHandler.HandleGRPCError(err, "role-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	r.logger.Info("gRPC client: successfully found role by name", zap.String("role_name", name))

	return r.mapper.ToApiResponseRole(res), nil
}

func (r *roleGrpcClientHandler) FindById(ctx context.Context, id int) (*response.ApiResponseRole, *response.ErrorResponse) {
	if errResp := r.errorHandler.ValidateClient(r.client, "role-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbrole.FindByIdRoleRequest{
		RoleId: int32(id),
	}

	res, err := r.client.FindByIdRole(ctx, req)

	if err != nil {
		return nil, r.errorHandler.HandleGRPCError(err, "role-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	r.logger.Info("gRPC client: successfully found role by id", zap.Int("role_id", id))

	return r.mapper.ToApiResponseRole(res), nil
}
