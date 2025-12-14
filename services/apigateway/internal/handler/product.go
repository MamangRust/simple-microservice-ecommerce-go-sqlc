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
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productHandleApi struct {
	queryClient    pb.ProductQueryServiceClient
	commandClient  pb.ProductCommandServiceClient
	logger         logger.LoggerInterface
	grpcmiddleware middlewares.GRPCErrorHandlingMiddleware
	mapping        mapper.ProductResponseMapper
	mencache       mencache.ProductCache
	observability  observability.TraceLoggerObservability
}

func NewProductHandle(router *echo.Echo, queryClient pb.ProductQueryServiceClient, commandClient pb.ProductCommandServiceClient, logger logger.LoggerInterface, grpcmiddleware middlewares.GRPCErrorHandlingMiddleware, mencache mencache.ProductCache) *productHandleApi {
	observability, _ := observability.NewObservability("product-service", logger)

	productHandler := &productHandleApi{
		queryClient:    queryClient,
		commandClient:  commandClient,
		logger:         logger,
		grpcmiddleware: grpcmiddleware,
		mapping:        mapper.NewProductResponseMapper(),
		mencache:       mencache,
		observability:  observability,
	}

	routercategory := router.Group("/api/product")

	routercategory.GET("", productHandler.FindAllProduct)
	routercategory.GET("/:id", productHandler.FindById)

	routercategory.GET("/active", productHandler.FindByActive)
	routercategory.GET("/trashed", productHandler.FindByTrashed)

	routercategory.POST("/create", productHandler.Create)
	routercategory.POST("/update/:id", productHandler.Update)

	routercategory.POST("/trashed/:id", productHandler.TrashedProduct)
	routercategory.POST("/restore/:id", productHandler.RestoreProduct)
	routercategory.DELETE("/permanent/:id", productHandler.DeleteProductPermanent)

	routercategory.POST("/restore/all", productHandler.RestoreAllProduct)
	routercategory.POST("/permanent/all", productHandler.DeleteAllProductPermanent)

	return productHandler
}

// @Security Bearer
// @Summary Find all products
// @Tags Product
// @Description Retrieve a list of all products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product [get]
func (h *productHandleApi) FindAllProduct(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "product"); err != nil {
		return err
	}

	const method = "FindAllProduct"

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

	logSuccess("Attempting to find all products",
		zap.String("handler", "FindAllProduct"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedProducts(ctx, reqCache); found {
		logSuccess("All products found in cache", zap.String("handler", "FindAllProduct"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindAll(ctx, req)
	if err != nil {
		logError("Failed to find all products via query client", err, zap.Error(err))
		return errors.ErrApiProductFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)

	h.mencache.SetCachedProducts(ctx, reqCache, so)

	logSuccess("Successfully found all products", zap.String("handler", "FindAllProduct"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find product by ID
// @Tags Product
// @Description Retrieve a product by ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProduct "Product data"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/{id} [get]
func (h *productHandleApi) FindById(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "product"); err != nil {
		return err
	}

	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid product ID provided in FindById", err, zap.String("handler", "FindById"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiProductInvalidId(c)
	}

	logSuccess("Attempting to find product by ID", zap.String("handler", "FindById"), zap.Int("product_id", id))

	if data, found := h.mencache.GetCachedProduct(ctx, id); found {
		logSuccess("Product found in cache", zap.String("handler", "FindById"), zap.Int("product_id", id))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindByIdProductRequest{Id: int32(id)}
	res, err := h.queryClient.FindById(ctx, req)
	if err != nil {
		logError("Failed to find product by ID via query client", err, zap.Int("product_id", id), zap.Error(err))
		return errors.ErrApiProductFailedFindById(c)
	}

	so := h.mapping.ToApiResponseProduct(res)

	h.mencache.SetCachedProduct(ctx, so)

	logSuccess("Successfully found product by ID", zap.String("handler", "FindById"), zap.Int("product_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active products
// @Tags Product
// @Description Retrieve a list of active products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProductDeleteAt "List of active products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/active [get]
func (h *productHandleApi) FindByActive(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "product"); err != nil {
		return err
	}
	ctx := c.Request().Context()
	const method = "FindByActive"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	logSuccess("Attempting to find active products",
		zap.String("handler", "FindByActive"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedProductActive(ctx, reqCache); found {
		logSuccess("Active products found in cache", zap.String("handler", "FindByActive"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByActive(ctx, req)
	if err != nil {
		logError("Failed to find active products via query client", err, zap.Error(err))
		return errors.ErrApiProductFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationProductDeleteAt(res)

	h.mencache.SetCachedProductActive(ctx, reqCache, so)

	logSuccess("Successfully found active products", zap.String("handler", "FindByActive"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed products
// @Tags Product
// @Description Retrieve a list of trashed products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProductDeleteAt "List of trashed products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/trashed [get]
func (h *productHandleApi) FindByTrashed(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "product"); err != nil {
		return err
	}
	ctx := c.Request().Context()
	const method = "FindByTrashed"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	logSuccess("Attempting to find trashed products",
		zap.String("handler", "FindByTrashed"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedProductTrashed(ctx, reqCache); found {
		logSuccess("Trashed products found in cache", zap.String("handler", "FindByTrashed"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByTrashed(ctx, req)
	if err != nil {
		logError("Failed to find trashed products via query client", err, zap.Error(err))
		return errors.ErrApiProductFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationProductDeleteAt(res)

	h.mencache.SetCachedProductTrashed(ctx, reqCache, so)

	logSuccess("Successfully found trashed products", zap.String("handler", "FindByTrashed"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// CreateProduct handles the creation of a new product with image upload.
// @Summary Create a new product
// @Tags Product
// @Description Create a new product with the provided details and an image file
// @Accept json
// @Produce json
// @Param request body requests.CreateProductRequest true "Product details"
// @Success 200 {object} response.ApiResponseProduct "Successfully created product"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create product"
// @Router /api/product/create [post]
func (h *productHandleApi) Create(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	ctx := c.Request().Context()

	const method = "Create"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	var req requests.CreateProductRequest

	if err := c.Bind(&req); err != nil {
		logError("Failed to bind CreateProduct request", err, zap.Error(err))
		return errors.ErrApiBindCreateProduct(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate CreateProduct request", err, zap.String("name", req.Name), zap.Error(err))
		return errors.ErrApiValidateCreateProduct(c)
	}

	logSuccess("Attempting to create a new product", zap.String("handler", "Create"), zap.String("name", req.Name), zap.Int64("price", int64(req.Price)), zap.Int("stock", req.Stock))

	grpcReq := &pb.CreateProductRequest{
		Name:  req.Name,
		Price: int64(req.Price),
		Stock: int32(req.Stock),
	}

	res, err := h.commandClient.Create(ctx, grpcReq)
	if err != nil {
		logError("Failed to create product via command client", err, zap.String("name", req.Name), zap.Error(err))
		return errors.ErrApiProductFailedCreate(c)
	}

	so := h.mapping.ToApiResponseProduct(res)
	logSuccess("Successfully created new product", zap.String("handler", "Create"), zap.String("name", req.Name), zap.Int32("new_product_id", int32(so.Data.ID)))

	return c.JSON(http.StatusCreated, so)
}

// @Security Bearer
// UpdateProduct handles the update of an existing product with optional image upload.
// @Summary Update an existing product
// @Tags Product
// @Description Update an existing product record with the provided details and an optional image file
// @Accept json
// @Produce json
// @Param request body requests.UpdateProductRequest true "Order details"
// @Success 200 {object} response.ApiResponseProduct "Successfully updated product"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update product"
// @Router /api/product/update/{id} [post]
func (h *productHandleApi) Update(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid product ID provided in Update", err, zap.String("handler", "Update"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiProductInvalidId(c)
	}

	var req requests.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		logError("Failed to bind UpdateProduct request", err, zap.Error(err))
		return errors.ErrApiBindUpdateProduct(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate UpdateProduct request", err, zap.Int("product_id", idInt), zap.String("name", req.Name), zap.Error(err))
		return errors.ErrApiValidateUpdateProduct(c)
	}

	logSuccess("Attempting to update product", zap.String("handler", "Update"), zap.Int("product_id", idInt), zap.String("name", req.Name))

	grpcReq := &pb.UpdateProductRequest{
		Id:    int32(idInt),
		Name:  req.Name,
		Price: int64(req.Price),
		Stock: int32(req.Stock),
	}

	res, err := h.commandClient.Update(ctx, grpcReq)
	if err != nil {
		logError("Failed to update product via command client", err, zap.Int("product_id", idInt), zap.Error(err))
		return errors.ErrApiProductFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseProduct(res)
	logSuccess("Successfully updated product", zap.String("handler", "Update"), zap.Int("product_id", idInt))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedProduct retrieves a trashed product record by its ID.
// @Summary Retrieve a trashed product
// @Tags Product
// @Description Retrieve a trashed product record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDeleteAt "Successfully retrieved trashed product"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed product"
// @Router /api/product/trashed/{id} [get]
func (h *productHandleApi) TrashedProduct(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	ctx := c.Request().Context()

	const method = "TrashedProduct"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid product ID provided in TrashedProduct", err, zap.String("handler", "TrashedProduct"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiProductInvalidId(c)
	}

	logSuccess("Attempting to trash product", zap.String("handler", "TrashedProduct"), zap.Int("product_id", id))

	req := &pb.FindByIdProductRequest{Id: int32(id)}
	res, err := h.commandClient.Trashed(ctx, req)
	if err != nil {
		logError("Failed to trash product via command client", err, zap.Int("product_id", id), zap.Error(err))
		return errors.ErrApiProductFailedTrashed(c)
	}

	so := h.mapping.ToApiResponsesProductDeleteAt(res)
	logSuccess("Successfully trashed product", zap.String("handler", "TrashedProduct"), zap.Int("product_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreProduct restores a product record from the trash by its ID.
// @Summary Restore a trashed product
// @Tags Product
// @Description Restore a trashed product record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDeleteAt "Successfully restored product"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore product"
// @Router /api/product/restore/{id} [post]
func (h *productHandleApi) RestoreProduct(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "RestoreProduct"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid product ID provided in RestoreProduct", err, zap.String("handler", "RestoreProduct"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiProductInvalidId(c)
	}

	logSuccess("Attempting to restore product", zap.String("handler", "RestoreProduct"), zap.Int("product_id", id))

	req := &pb.FindByIdProductRequest{Id: int32(id)}
	res, err := h.commandClient.Restore(ctx, req)
	if err != nil {
		logError("Failed to restore product via command client", err, zap.Int("product_id", id), zap.Error(err))
		return errors.ErrApiProductFailedRestore(c)
	}

	so := h.mapping.ToApiResponsesProductDeleteAt(res)
	logSuccess("Successfully restored product", zap.String("handler", "RestoreProduct"), zap.Int("product_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteProductPermanent permanently deletes a product record by its ID.
// @Summary Permanently delete a product
// @Tags Product
// @Description Permanently delete a product record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDelete "Successfully deleted product record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete product:"
// @Router /api/product/delete/{id} [delete]
func (h *productHandleApi) DeleteProductPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "DeleteProductPermanent"
	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid product ID provided in DeleteProductPermanent", err, zap.String("handler", "DeleteProductPermanent"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiProductInvalidId(c)
	}

	logSuccess("Attempting to permanently delete product", zap.String("handler", "DeleteProductPermanent"), zap.Int("product_id", id))

	req := &pb.FindByIdProductRequest{Id: int32(id)}
	res, err := h.commandClient.DeleteProductPermanent(ctx, req)
	if err != nil {
		logError("Failed to permanently delete product via command client", err, zap.Int("product_id", id), zap.Error(err))
		return errors.ErrApiProductFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseProductDelete(res)
	logSuccess("Successfully permanently deleted product", zap.String("handler", "DeleteProductPermanent"), zap.Int("product_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllProduct restores all trashed product records.
// @Summary Restore all trashed products
// @Tags Product
// @Description Restore all trashed product records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseProductAll "Successfully restored all products"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all products"
// @Router /api/product/restore/all [post]
func (h *productHandleApi) RestoreAllProduct(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "RestoreAllProduct"
	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to restore all trashed products", zap.String("handler", "RestoreAllProduct"))

	res, err := h.commandClient.RestoreAllProduct(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to restore all products via command client", err, zap.Error(err))
		return errors.ErrApiProductFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseProductAll(res)
	logSuccess("Successfully restored all trashed products", zap.String("handler", "RestoreAllProduct"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllProductPermanent permanently deletes all product records.
// @Summary Permanently delete all products
// @Tags Product
// @Description Permanently delete all product records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseProductAll "Successfully deleted all product records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all products"
// @Router /api/product/delete/all [post]
func (h *productHandleApi) DeleteAllProductPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "product"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "DeleteAllProductPermanent"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to permanently delete all trashed products", zap.String("handler", "DeleteAllProductPermanent"))

	res, err := h.commandClient.DeleteAllProduct(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to permanently delete all products via command client", err, zap.Error(err))
		return errors.ErrApiProductFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseProductAll(res)
	logSuccess("Successfully permanently deleted all trashed products", zap.String("handler", "DeleteAllProductPermanent"))

	return c.JSON(http.StatusOK, so)
}
