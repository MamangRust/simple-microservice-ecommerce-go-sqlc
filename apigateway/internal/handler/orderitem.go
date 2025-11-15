package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/mapper"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/errors"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type orderItemHandleApi struct {
	client         pb.OrderItemServiceClient
	logger         logger.LoggerInterface
	grpcmiddleware middlewares.GRPCErrorHandlingMiddleware
	mapping        mapper.OrderItemResponseMapper
}

func NewOrderItemHandle(router *echo.Echo, client pb.OrderItemServiceClient, logger logger.LoggerInterface, grpcmiddleware middlewares.GRPCErrorHandlingMiddleware) *orderItemHandleApi {
	orderHandle := &orderItemHandleApi{
		client:         client,
		logger:         logger,
		grpcmiddleware: grpcmiddleware,
		mapping:        mapper.NewOrderItemResponseMapper(),
	}

	routercategory := router.Group("/api/order-item")

	routercategory.GET("", orderHandle.FindAllOrderItems)
	routercategory.GET("/:order_id", orderHandle.FindOrderItemByOrder)
	routercategory.GET("/active", orderHandle.FindByActive)
	routercategory.GET("/trashed", orderHandle.FindByTrashed)

	return orderHandle

}

// @Security Bearer
// @Summary Find all order items
// @Tags OrderItem
// @Description Retrieve a list of all order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItem "List of order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item [get]
func (h *orderItemHandleApi) FindAllOrderItems(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "order_item"); err != nil {
		return err
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	h.logger.Info("Attempting to find all order items",
		zap.String("handler", "FindAllOrderItems"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	ctx := c.Request().Context()
	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		h.logger.Error("Failed to find all order items via client", zap.Error(err))
		return errors.ErrApiOrderItemFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItem(res)
	h.logger.Info("Successfully found all order items", zap.String("handler", "FindAllOrderItems"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active order items
// @Tags OrderItem
// @Description Retrieve a list of active order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of active order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/active [get]
func (h *orderItemHandleApi) FindByActive(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "order_item"); err != nil {
		return err
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	h.logger.Info("Attempting to find active order items",
		zap.String("handler", "FindByActive"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	ctx := c.Request().Context()
	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		h.logger.Error("Failed to find active order items via client", zap.Error(err))
		return errors.ErrApiOrderItemFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItemDeleteAt(res)
	h.logger.Info("Successfully found active order items", zap.String("handler", "FindByActive"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed order items
// @Tags OrderItem
// @Description Retrieve a list of trashed order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of trashed order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/trashed [get]
func (h *orderItemHandleApi) FindByTrashed(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "order_item"); err != nil {
		return err
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	h.logger.Info("Attempting to find trashed order items",
		zap.String("handler", "FindByTrashed"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	ctx := c.Request().Context()
	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		h.logger.Error("Failed to find trashed order items via client", zap.Error(err))
		return errors.ErrApiOrderItemFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItemDeleteAt(res)
	h.logger.Info("Successfully found trashed order items", zap.String("handler", "FindByTrashed"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find order items by order ID
// @Tags OrderItem
// @Description Retrieve order items by order ID
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponsesOrderItem "List of order items by order ID"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/order/{order_id} [get]
func (h *orderItemHandleApi) FindOrderItemByOrder(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "order_item"); err != nil {
		return err
	}

	ctx := c.Request().Context()

	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		h.logger.Warn("Invalid order ID provided in FindOrderItemByOrder", zap.String("handler", "FindOrderItemByOrder"), zap.String("param_order_id", c.Param("order_id")))
		return errors.ErrApiOrderItemInvalidId(c)
	}

	h.logger.Info("Attempting to find order items by order ID", zap.String("handler", "FindOrderItemByOrder"), zap.Int("order_id", orderID))

	req := &pb.FindByIdOrderItemRequest{
		Id: int32(orderID),
	}

	res, err := h.client.FindOrderItemByOrder(ctx, req)
	if err != nil {
		h.logger.Error("Failed to find order items by order ID via client", zap.Int("order_id", orderID), zap.Error(err))
		return errors.ErrApiOrderItemFailedFindByOrderId(c)
	}

	so := h.mapping.ToApiResponsesOrderItem(res)
	h.logger.Info("Successfully found order items by order ID", zap.String("handler", "FindOrderItemByOrder"), zap.Int("order_id", orderID))

	return c.JSON(http.StatusOK, so)
}
