package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	grpcclientmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/grpcclient"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/middlewares"

	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	"go.uber.org/zap"
)

type ProductGrpcClientHandler interface {
	FindById(ctx context.Context, id int32) (*response.ApiResponseProduct, *response.ErrorResponse)
	UpdateProductStock(ctx context.Context, req *requests.UpdateProductStockRequest) (*response.ApiResponseProduct, *response.ErrorResponse)
}

type productGrpcClientHandler struct {
	clientQuery   pbproduct.ProductQueryServiceClient
	clientCommand pbproduct.ProductCommandServiceClient
	logger        logger.LoggerInterface
	errorHandler  middlewares.GRPCErrorHandling
	mapper        grpcclientmapper.ProductClientResponseMapper
}

func NewProductGrpcClientHandler(clientQuery pbproduct.ProductQueryServiceClient, clientCommand pbproduct.ProductCommandServiceClient, logger logger.LoggerInterface, errorHandler middlewares.GRPCErrorHandling) ProductGrpcClientHandler {
	mapper := grpcclientmapper.NewProductClientResponseMapper()

	return &productGrpcClientHandler{
		clientQuery:   clientQuery,
		clientCommand: clientCommand,
		logger:        logger,
		errorHandler:  errorHandler,
		mapper:        mapper,
	}
}

func (h *productGrpcClientHandler) FindById(ctx context.Context, id int32) (*response.ApiResponseProduct, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "product-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbproduct.FindByIdProductRequest{Id: id}

	res, err := h.clientQuery.FindById(ctx, req)
	if err != nil {
		return nil, h.errorHandler.HandleGRPCError(err, "product-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	h.logger.Info("gRPC client: successfully found product", zap.Int32("product_id", id))

	return h.mapper.ToApiResponseProduct(res), nil
}

func (h productGrpcClientHandler) UpdateProductStock(ctx context.Context, req *requests.UpdateProductStockRequest) (*response.ApiResponseProduct, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientCommand, "product-service"); errResp != nil {
		return nil, errResp
	}

	reqPb := &pbproduct.UpdateProductStockRequest{
		Id:    int32(req.ProductID),
		Stock: int32(req.Stock),
	}

	res, err := h.clientCommand.UpdateProductCountStock(ctx, reqPb)

	if err != nil {
		return nil, h.errorHandler.HandleGRPCError(err, "product-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	h.logger.Info("gRPC client: successfully updated product stock",
		zap.Int("product_id", req.ProductID),
		zap.Int32("updated_stock", res.Data.Stock),
	)

	return h.mapper.ToApiResponseProduct(res), nil
}
