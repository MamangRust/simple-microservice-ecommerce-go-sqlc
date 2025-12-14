package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/mapper"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/middlewares"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/errors"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/observability"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleApi struct {
	queryClient    pb.OrderQueryServiceClient
	commandClient  pb.OrderCommandServiceClient
	logger         logger.LoggerInterface
	grpcmiddleware middlewares.GRPCErrorHandlingMiddleware
	mapping        mapper.OrderResponseMapper
	mencache       mencache.OrderCache
	observability  observability.TraceLoggerObservability
}

func NewOrderHandle(router *echo.Echo, queryClient pb.OrderQueryServiceClient, commandClient pb.OrderCommandServiceClient, logger logger.LoggerInterface, grpcmiddleware middlewares.GRPCErrorHandlingMiddleware, mencache mencache.OrderCache) *orderHandleApi {
	observability, _ := observability.NewObservability("order-service", logger)

	orderHandler := &orderHandleApi{
		queryClient:    queryClient,
		commandClient:  commandClient,
		logger:         logger,
		grpcmiddleware: grpcmiddleware,
		mapping:        mapper.NewOrderResponseMapper(),
		mencache:       mencache,
		observability:  observability,
	}

	routerOrder := router.Group("/api/order")

	routerOrder.GET("", orderHandler.FindAllOrders)
	routerOrder.GET("/:id", orderHandler.FindById)
	routerOrder.GET("/active", orderHandler.FindByActive)
	routerOrder.GET("/trashed", orderHandler.FindByTrashed)

	routerOrder.POST("/create", orderHandler.Create)
	routerOrder.POST("/update/:id", orderHandler.Update)

	routerOrder.POST("/trashed/:id", orderHandler.TrashedOrder)
	routerOrder.POST("/restore/:id", orderHandler.RestoreOrder)
	routerOrder.DELETE("/permanent/:id", orderHandler.DeleteOrderPermanent)

	routerOrder.POST("/restore/all", orderHandler.RestoreAllOrder)
	routerOrder.POST("/permanent/all", orderHandler.DeleteAllOrderPermanent)

	return orderHandler
}

// @Security Bearer
// @Summary Find all orders
// @Tags Order
// @Description Retrieve a list of all orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrder "List of orders"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order [get]
func (h *orderHandleApi) FindAllOrders(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "order"); err != nil {
		return err
	}
	const method = "FindAllOrders"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	logSuccess("Attempting to find all orders",
		zap.String("handler", "FindAll"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetOrderAllCache(ctx, reqCache); found {
		logSuccess("All orders found in cache", zap.String("handler", "FindAll"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindAll(ctx, req)
	if err != nil {
		logError("Failed to find all orders via query client", err, zap.Error(err))
		return errors.ErrApiOrderFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationOrder(res)

	h.mencache.SetOrderAllCache(ctx, reqCache, so)

	logSuccess("Successfully found all orders", zap.String("handler", "FindAll"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find order by ID
// @Tags Order
// @Description Retrieve an order by ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrder "Order data"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order/{id} [get]
func (h *orderHandleApi) FindById(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "order"); err != nil {
		return err
	}

	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid order ID provided in FindById", err, zap.String("handler", "FindById"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiOrderInvalidId(c)
	}

	logSuccess("Attempting to find order by ID", zap.String("handler", "FindById"), zap.Int("order_id", id))

	if data, found := h.mencache.GetCachedOrderCache(ctx, id); found {
		logSuccess("Order found in cache", zap.String("handler", "FindById"), zap.Int("order_id", id))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindByIdOrderRequest{Id: int32(id)}
	res, err := h.queryClient.FindById(ctx, req)
	if err != nil {
		logError("Failed to find order by ID via query client", err, zap.Int("order_id", id), zap.Error(err))
		return errors.ErrApiOrderFailedFindById(c)
	}

	so := h.mapping.ToApiResponseOrder(res)

	h.mencache.SetCachedOrderCache(ctx, so)

	logSuccess("Successfully found order by ID", zap.String("handler", "FindById"), zap.Int("order_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active orders
// @Tags Order
// @Description Retrieve a list of active orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt "List of active orders"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order/active [get]
func (h *orderHandleApi) FindByActive(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "order"); err != nil {
		return err
	}
	const method = "FindByActive"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	logSuccess("Attempting to find active orders",
		zap.String("handler", "FindByActive"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetOrderActiveCache(ctx, reqCache); found {
		logSuccess("Active orders found in cache", zap.String("handler", "FindByActive"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByActive(ctx, req)
	if err != nil {
		logError("Failed to find active orders via query client", err, zap.Error(err))
		return errors.ErrApiOrderFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	h.mencache.SetOrderActiveCache(ctx, reqCache, so)

	logSuccess("Successfully found active orders", zap.String("handler", "FindByActive"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed orders
// @Tags Order
// @Description Retrieve a list of trashed orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt "List of trashed orders"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order data"
// @Router /api/order/trashed [get]
func (h *orderHandleApi) FindByTrashed(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "order"); err != nil {
		return err
	}
	const method = "FindByTrashed"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	logSuccess("Attempting to find trashed orders",
		zap.String("handler", "FindByTrashed"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetOrderTrashedCache(ctx, reqCache); found {
		logSuccess("Trashed orders found in cache", zap.String("handler", "FindByTrashed"))
		return c.JSON(http.StatusOK, data)
	}
	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByTrashed(ctx, req)
	if err != nil {
		logError("Failed to find trashed orders via query client", err, zap.Error(err))
		return errors.ErrApiOrderFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	h.mencache.SetOrderTrashedCache(ctx, reqCache, so)

	logSuccess("Successfully found trashed orders", zap.String("handler", "FindByTrashed"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new order
// @Tags Order
// @Description Create a new order with provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateOrderRequest true "Order details"
// @Success 200 {object} response.ApiResponseOrder "Successfully created order"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create order"
// @Router /api/order/create [post]
func (h *orderHandleApi) Create(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	var req requests.CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		logError("Failed to bind CreateOrder request", err, zap.Error(err))
		return errors.ErrApiBindCreateOrder(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate CreateOrder request", err, zap.Int("user_id", req.UserID), zap.Error(err))
		return errors.ErrApiValidateCreateOrder(c)
	}

	logSuccess("Attempting to create a new order", zap.String("handler", "Create"), zap.Int("user_id", req.UserID))

	grpcReq := &pb.CreateOrderRequest{
		UserId: int32(req.UserID),
		Items:  make([]*pb.CreateOrderItemRequest, len(req.Items)),
	}

	for i, item := range req.Items {
		grpcReq.Items[i] = &pb.CreateOrderItemRequest{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
		}
	}

	res, err := h.commandClient.Create(ctx, grpcReq)
	if err != nil {
		logError("Failed to create order via command client", err, zap.Int("user_id", req.UserID), zap.Error(err))
		return errors.ErrApiOrderFailedCreate(c)
	}

	so := h.mapping.ToApiResponseOrder(res)
	logSuccess("Successfully created new order", zap.String("handler", "Create"), zap.Int32("new_order_id", int32(so.Data.ID)))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing order
// @Tags Order
// @Description Update an existing order with provided details
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body requests.UpdateOrderRequest true "Order update details"
// @Success 200 {object} response.ApiResponseOrder "Successfully updated order"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update order"
// @Router /api/order/update [post]
func (h *orderHandleApi) Update(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}

	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid order ID provided in Update", err, zap.String("handler", "Update"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiOrderInvalidId(c)
	}

	var req requests.UpdateOrderRequest
	req.OrderID = &idInt

	if err := c.Bind(&req); err != nil {
		logError("Failed to bind UpdateOrder request", err, zap.Error(err))
		return errors.ErrApiBindUpdateOrder(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate UpdateOrder request", err, zap.Int("order_id", idInt), zap.Error(err))
		return errors.ErrApiValidateUpdateOrder(c)
	}

	logSuccess("Attempting to update order", zap.String("handler", "Update"), zap.Int("order_id", idInt), zap.Int("user_id", req.UserID))

	grpcReq := &pb.UpdateOrderRequest{
		OrderId: int32(idInt),
		UserId:  int32(req.UserID),
		Items:   make([]*pb.UpdateOrderItemRequest, len(req.Items)),
	}

	for i, item := range req.Items {
		grpcReq.Items[i] = &pb.UpdateOrderItemRequest{
			OrderItemId: int32(item.OrderItemID),
			ProductId:   int32(item.ProductID),
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		}
	}

	res, err := h.commandClient.Update(ctx, grpcReq)
	if err != nil {
		logError("Failed to update order via command client", err, zap.Int("order_id", idInt), zap.Error(err))
		return errors.ErrApiOrderFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseOrder(res)
	logSuccess("Successfully updated order", zap.String("handler", "Update"), zap.Int("order_id", idInt))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedOrder retrieves a trashed order record by its ID.
// @Summary Retrieve a trashed order
// @Tags Order
// @Description Retrieve a trashed order record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDeleteAt "Successfully retrieved trashed order"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed order"
// @Router /api/order/trashed/{id} [post]
func (h *orderHandleApi) TrashedOrder(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}

	const method = "TrashedOrder"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid order ID provided in TrashedOrder", err, zap.String("handler", "TrashedOrder"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiOrderInvalidId(c)
	}

	logSuccess("Attempting to trash order", zap.String("handler", "TrashedOrder"), zap.Int("order_id", id))

	req := &pb.FindByIdOrderRequest{Id: int32(id)}
	res, err := h.commandClient.Trashed(ctx, req)
	if err != nil {
		logError("Failed to trash order via command client", err, zap.Int("order_id", id), zap.Error(err))
		return errors.ErrApiOrderFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)
	logSuccess("Successfully trashed order", zap.String("handler", "TrashedOrder"), zap.Int("order_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreOrder restores an order record from the trash by its ID.
// @Summary Restore a trashed order
// @Tags Order
// @Description Restore a trashed order record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDeleteAt "Successfully restored order"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore order"
// @Router /api/order/restore/{id} [post]
func (h *orderHandleApi) RestoreOrder(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}
	const method = "RestoreOrder"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid order ID provided in RestoreOrder", err, zap.String("handler", "RestoreOrder"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiOrderInvalidId(c)
	}

	logSuccess("Attempting to restore order", zap.String("handler", "RestoreOrder"), zap.Int("order_id", id))

	req := &pb.FindByIdOrderRequest{Id: int32(id)}
	res, err := h.commandClient.Restore(ctx, req)
	if err != nil {
		logError("Failed to restore order via command client", err, zap.Int("order_id", id), zap.Error(err))
		return errors.ErrApiOrderFailedRestore(c)
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)
	logSuccess("Successfully restored order", zap.String("handler", "RestoreOrder"), zap.Int("order_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteOrderPermanent permanently deletes an order record by its ID.
// @Summary Permanently delete an order
// @Tags Order
// @Description Permanently delete an order record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDelete "Successfully deleted order record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete order:"
// @Router /api/order/delete/{id} [delete]
func (h *orderHandleApi) DeleteOrderPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}

	const method = "DeleteOrderPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid order ID provided in DeleteOrderPermanent", err, zap.String("handler", "DeleteOrderPermanent"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiOrderInvalidId(c)
	}

	logSuccess("Attempting to permanently delete order", zap.String("handler", "DeleteOrderPermanent"), zap.Int("order_id", id))

	req := &pb.FindByIdOrderRequest{Id: int32(id)}
	res, err := h.commandClient.DeleteOrderPermanent(ctx, req)
	if err != nil {
		logError("Failed to permanently delete order via command client", err, zap.Int("order_id", id), zap.Error(err))
		return errors.ErrApiOrderFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseOrderDelete(res)
	logSuccess("Successfully permanently deleted order", zap.String("handler", "DeleteOrderPermanent"), zap.Int("order_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllOrder restores all trashed orders.
// @Summary Restore all trashed orders
// @Tags Order
// @Description Restore all trashed order records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll "Successfully restored all orders"
// @Failure 500 {object} response.ErrorResponse "Failed to restore orders"
// @Router /api/order/restore/all [post]
func (h *orderHandleApi) RestoreAllOrder(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}

	const method = "RestoreAllOrder"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to restore all trashed orders", zap.String("handler", "RestoreAllOrder"))

	res, err := h.commandClient.RestoreAllOrder(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to restore all orders via command client", err, zap.Error(err))
		return errors.ErrApiOrderFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseOrderAll(res)
	logSuccess("Successfully restored all trashed orders", zap.String("handler", "RestoreAllOrder"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllOrderPermanent permanently deletes all orders.
// @Summary Permanently delete all orders
// @Tags Order
// @Description Permanently delete all order records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll "Successfully deleted all orders permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete orders"
// @Router /api/order/delete/all [post]
func (h *orderHandleApi) DeleteAllOrderPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "order"); err != nil {
		return err
	}

	const method = "DeleteAllOrderPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to permanently delete all trashed orders", zap.String("handler", "DeleteAllOrderPermanent"))

	res, err := h.commandClient.DeleteAllOrder(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to permanently delete all orders via command client", err, zap.Error(err))
		return errors.ErrApiOrderFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseOrderAll(res)
	logSuccess("Successfully permanently deleted all trashed orders", zap.String("handler", "DeleteAllOrderPermanent"))

	return c.JSON(http.StatusOK, so)
}
