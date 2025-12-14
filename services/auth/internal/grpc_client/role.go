package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpclientmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/response/grpclient"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
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
		mapper:       mapper,
		errorHandler: errorHandler,
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
		r.logger.Error("gRPC client error when finding role by name", zap.String("name", name), zap.Error(err))

		return nil, r.errorHandler.HandleGRPCError(err, "role-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return r.mapper.ToApiResponseRole(res), nil
}
