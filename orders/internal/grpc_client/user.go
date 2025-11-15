package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	grpcclientmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/grpcclient"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"go.uber.org/zap"
)

type UserGrpcClientHandler interface {
	FindById(ctx context.Context, id int32) (*response.ApiResponseUser, *response.ErrorResponse)
}

type userGrpcClientHandler struct {
	client       pbuser.UserQueryServiceClient
	logger       logger.LoggerInterface
	errorHandler middlewares.GRPCErrorHandling
	mapper       grpcclientmapper.UserClientResponseMapper
}

func NewUserGrpcClientHandler(userClient pbuser.UserQueryServiceClient, logger logger.LoggerInterface, errorHandler middlewares.GRPCErrorHandling) UserGrpcClientHandler {
	mapper := grpcclientmapper.NewUserClientResponseMapper()

	return &userGrpcClientHandler{
		client:       userClient,
		logger:       logger,
		errorHandler: errorHandler,
		mapper:       mapper,
	}
}

func (h *userGrpcClientHandler) FindById(ctx context.Context, id int32) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.client, "user-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuser.FindByIdUserRequest{Id: id}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	h.logger.Info("gRPC client: successfully found user", zap.Int32("user_id", id))

	return h.mapper.ToApiResponseUser(res), nil
}
