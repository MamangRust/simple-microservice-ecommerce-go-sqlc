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
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleApi struct {
	queryClient    pb.RoleQueryServiceClient
	commandClient  pb.RoleCommandServiceClient
	logger         logger.LoggerInterface
	grpcmiddleware middlewares.GRPCErrorHandlingMiddleware
	mapping        mapper.RoleResponseMapper
	mencache       mencache.RoleCache
	observability  observability.TraceLoggerObservability
}

func NewRoleHandleApi(router *echo.Echo, queryClient pb.RoleQueryServiceClient,
	commandClient pb.RoleCommandServiceClient, logger logger.LoggerInterface, grpcmiddleware middlewares.GRPCErrorHandlingMiddleware, mencache mencache.RoleCache) *roleHandleApi {
	observability, _ := observability.NewObservability("role-service", logger)

	roleHandler := &roleHandleApi{
		queryClient:    queryClient,
		commandClient:  commandClient,
		logger:         logger,
		grpcmiddleware: grpcmiddleware,
		mapping:        mapper.NewRoleResponseMapper(),
		mencache:       mencache,
		observability:  observability,
	}

	routerRole := router.Group("/api/role")

	routerRole.GET("", roleHandler.FindAll)
	routerRole.GET("/:id", roleHandler.FindById)
	routerRole.GET("/active", roleHandler.FindByActive)
	routerRole.GET("/trashed", roleHandler.FindByTrashed)
	routerRole.GET("/user/:user_id", roleHandler.FindByUserId)
	routerRole.POST("", roleHandler.Create)
	routerRole.POST("/:id", roleHandler.Update)
	routerRole.DELETE("/:id", roleHandler.Trashed)
	routerRole.PUT("/restore/:id", roleHandler.Restore)
	routerRole.DELETE("/permanent/:id", roleHandler.DeletePermanent)

	routerRole.POST("/restore/all", roleHandler.RestoreAll)
	routerRole.POST("/permanent/all", roleHandler.DeleteAllPermanent)

	return roleHandler
}

// FindAll godoc.
// @Summary Get all roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRole "List of roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch roles"
// @Router /api/role [get]
func (h *roleHandleApi) FindAll(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "role"); err != nil {
		return err
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
	)
	ctx := c.Request().Context()

	const method = "FindAll"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	logSuccess("Attempting to find all roles",
		zap.String("handler", "FindAll"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedRoles(ctx, reqCache); found {
		logSuccess("All roles found in cache", zap.String("handler", "FindAll"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindAllRole(ctx, req)
	if err != nil {
		logError("Failed to find all roles via query client", err, zap.Error(err))
		return errors.ErrApiFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationRole(res)

	h.mencache.SetCachedRoles(ctx, reqCache, so)

	logSuccess("Successfully found all roles", zap.String("handler", "FindAll"))

	return c.JSON(http.StatusOK, so)
}

// FindById godoc.
// @Summary Get a role by ID
// @Tags Role
// @Security Bearer
// @Description Retrieve a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch role"
// @Router /api/role/{id} [get]
func (h *roleHandleApi) FindById(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "FindById"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		logError("Invalid role ID provided in FindById", err, zap.String("handler", "FindById"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiRoleInvalidId(c)
	}

	logSuccess("Attempting to find role by ID", zap.String("handler", "FindById"), zap.Int("role_id", roleID))

	if data, found := h.mencache.GetCachedRoleById(ctx, roleID); found {
		logSuccess("Role found in cache", zap.String("handler", "FindById"), zap.Int("role_id", roleID))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.queryClient.FindByIdRole(ctx, req)
	if err != nil {
		logError("Failed to find role by ID via query client", err, zap.Int("role_id", roleID), zap.Error(err))
		return errors.ErrApiRoleNotFound(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	h.mencache.SetCachedRoleById(ctx, so)

	logSuccess("Successfully found role by ID", zap.String("handler", "FindById"), zap.Int("role_id", roleID))

	return c.JSON(http.StatusOK, so)
}

// FindByActive godoc.
// @Summary Get active roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of active roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch active roles"
// @Router /api/role/active [get]
func (h *roleHandleApi) FindByActive(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "role"); err != nil {
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

	logSuccess("Attempting to find active roles",
		zap.String("handler", "FindByActive"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedRoleActive(ctx, reqCache); found {
		logSuccess("Active roles found in cache", zap.String("handler", "FindByActive"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByActive(ctx, req)
	if err != nil {
		logError("Failed to find active roles via query client", err, zap.Error(err))
		return errors.ErrApiFailedFindActive(c)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	h.mencache.SetCachedRoleActive(ctx, reqCache, so)

	logSuccess("Successfully found active roles", zap.String("handler", "FindByActive"))

	return c.JSON(http.StatusOK, so)
}

// FindByTrashed godoc.
// @Summary Get trashed roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of trashed roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch trashed roles"
// @Router /api/role/trashed [get]
func (h *roleHandleApi) FindByTrashed(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "role"); err != nil {
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

	logSuccess("Attempting to find trashed roles",
		zap.String("handler", "FindByTrashed"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedRoleTrashed(ctx, reqCache); found {
		logSuccess("Trashed roles found in cache", zap.String("handler", "FindByTrashed"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByTrashed(ctx, req)
	if err != nil {
		logError("Failed to find trashed roles via query client", err, zap.Error(err))
		return errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	h.mencache.SetCachedRoleTrashed(ctx, reqCache, so)

	logSuccess("Successfully found trashed roles", zap.String("handler", "FindByTrashed"))

	return c.JSON(http.StatusOK, so)
}

// FindByUserId godoc.
// @Summary Get role by user ID
// @Tags Role
// @Security Bearer
// @Description Retrieve a role by the associated user ID.
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch role by user ID"
// @Router /api/role/user/{user_id} [get]
func (h *roleHandleApi) FindByUserId(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "FindByUserId"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		logError("Invalid user ID provided in FindByUserId", err, zap.String("handler", "FindByUserId"), zap.String("param_user_id", c.Param("user_id")))
		return errors.ErrApiRoleInvalidId(c)
	}

	logSuccess("Attempting to find roles by user ID", zap.String("handler", "FindByUserId"), zap.Int("user_id", userID))

	if data, found := h.mencache.GetCachedRoleByUserId(ctx, userID); found {
		logSuccess("Roles by user ID found in cache", zap.String("handler", "FindByUserId"), zap.Int("user_id", userID))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindByIdUserRoleRequest{
		UserId: int32(userID),
	}

	res, err := h.queryClient.FindByUserId(ctx, req)
	if err != nil {
		logError("Failed to find roles by user ID via query client", err, zap.Int("user_id", userID), zap.Error(err))
		return errors.ErrApiRoleNotFound(c)
	}

	so := h.mapping.ToApiResponsesRole(res)

	h.mencache.SetCachedRoleByUserId(ctx, userID, so)

	logSuccess("Successfully found roles by user ID", zap.String("handler", "FindByUserId"), zap.Int("user_id", userID))

	return c.JSON(http.StatusOK, so)
}

// Create godoc.
// @Summary Create a new role
// @Tags Role
// @Security Bearer
// @Description Create a new role with the provided details.
// @Accept json
// @Produce json
// @Param request body requests.CreateRoleRequest true "Role data"
// @Success 200 {object} response.ApiResponseRole "Created role data"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to create role"
// @Router /api/role/create [post]
func (h *roleHandleApi) Create(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "Create"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	var req requests.CreateRoleRequest

	if err := c.Bind(&req); err != nil {
		logError("Failed to bind CreateRole request", err, zap.Error(err))
		return errors.ErrApiBindCreateRole(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate CreateRole request", err, zap.String("name", req.Name), zap.Error(err))
		return errors.ErrApiValidateCreateRole(c)
	}

	logSuccess("Attempting to create a new role", zap.String("handler", "Create"), zap.String("name", req.Name))

	reqPb := &pb.CreateRoleRequest{
		Name: req.Name,
	}

	res, err := h.commandClient.CreateRole(ctx, reqPb)
	if err != nil {
		logError("Failed to create role via command client", err, zap.String("name", req.Name), zap.Error(err))
		return errors.ErrApiFailedCreateRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)
	logSuccess("Successfully created new role", zap.String("handler", "Create"), zap.String("name", req.Name))

	return c.JSON(http.StatusOK, so)
}

// Update godoc.
// @Summary Update a role
// @Tags Role
// @Security Bearer
// @Description Update an existing role with the provided details.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body requests.UpdateRoleRequest true "Role data"
// @Success 200 {object} response.ApiResponseRole "Updated role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID or request body"
// @Failure 500 {object} response.ErrorResponse "Failed to update role"
// @Router /api/role/update/{id} [post]
func (h *roleHandleApi) Update(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "Update"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		logError("Invalid role ID provided in Update", err, zap.String("handler", "Update"), zap.String("param_id", c.Param("id")))
		return errors.ErrInvalidRoleId(c)
	}

	var req requests.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		logError("Failed to bind UpdateRole request", err, zap.Error(err))
		return errors.ErrApiBindUpdateRole(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate UpdateRole request", err, zap.Int("role_id", roleID), zap.String("name", req.Name), zap.Error(err))
		return errors.ErrApiValidateUpdateRole(c)
	}

	logSuccess("Attempting to update role", zap.String("handler", "Update"), zap.Int("role_id", roleID), zap.String("name", req.Name))

	reqPb := &pb.UpdateRoleRequest{
		Id:   int32(roleID),
		Name: req.Name,
	}

	res, err := h.commandClient.UpdateRole(ctx, reqPb)
	if err != nil {
		logError("Failed to update role via command client", err, zap.Int("role_id", roleID), zap.Error(err))
		return errors.ErrApiFailedUpdateRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)
	logSuccess("Successfully updated role", zap.String("handler", "Update"), zap.Int("role_id", roleID))

	return c.JSON(http.StatusOK, so)
}

// Trashed godoc.
// @Summary Soft-delete a role
// @Tags Role
// @Security Bearer
// @Description Soft-delete a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Soft-deleted role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to soft-delete role"
// @Router /api/role/trashed/{id} [post]
func (h *roleHandleApi) Trashed(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "Trashed"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		logError("Invalid role ID provided in Trashed", err, zap.String("handler", "Trashed"), zap.String("param_id", c.Param("id")))
		return errors.ErrInvalidRoleId(c)
	}

	logSuccess("Attempting to trash role", zap.String("handler", "Trashed"), zap.Int("role_id", roleID))

	req := &pb.FindByIdRoleRequest{RoleId: int32(roleID)}
	res, err := h.commandClient.TrashedRole(ctx, req)
	if err != nil {
		logError("Failed to trash role via command client", err, zap.Int("role_id", roleID), zap.Error(err))
		return errors.ErrApiFailedTrashedRole(c)
	}

	so := h.mapping.ToApiResponseRoleDeleteAt(res)
	logSuccess("Successfully trashed role", zap.String("handler", "Trashed"), zap.Int("role_id", roleID))

	return c.JSON(http.StatusOK, so)
}

// Restore godoc.
// @Summary Restore a soft-deleted role
// @Tags Role
// @Security Bearer
// @Description Restore a soft-deleted role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Restored role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore role"
// @Router /api/role/restore/{id} [post]
func (h *roleHandleApi) Restore(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "Restore"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		logError("Invalid role ID provided in Restore", err, zap.String("handler", "Restore"), zap.String("param_id", c.Param("id")))
		return errors.ErrInvalidRoleId(c)
	}

	logSuccess("Attempting to restore role", zap.String("handler", "Restore"), zap.Int("role_id", roleID))

	req := &pb.FindByIdRoleRequest{RoleId: int32(roleID)}
	res, err := h.commandClient.RestoreRole(ctx, req)
	if err != nil {
		logError("Failed to restore role via command client", err, zap.Int("role_id", roleID), zap.Error(err))
		return errors.ErrApiFailedRestoreRole(c)
	}

	so := h.mapping.ToApiResponseRoleDeleteAt(res)
	logSuccess("Successfully restored role", zap.String("handler", "Restore"), zap.Int("role_id", roleID))

	return c.JSON(http.StatusOK, so)
}

// DeletePermanent godoc.
// @Summary Permanently delete a role
// @Tags Role
// @Security Bearer
// @Description Permanently delete a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Permanently deleted role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete role permanently"
// @Router /api/role/permanent/{id} [delete]
func (h *roleHandleApi) DeletePermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}
	const method = "DeletePermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		logError("Invalid role ID provided in DeletePermanent", err, zap.String("handler", "DeletePermanent"), zap.String("param_id", c.Param("id")))
		return errors.ErrInvalidRoleId(c)
	}

	logSuccess("Attempting to permanently delete role", zap.String("handler", "DeletePermanent"), zap.Int("role_id", roleID))

	req := &pb.FindByIdRoleRequest{RoleId: int32(roleID)}
	res, err := h.commandClient.DeleteRolePermanent(ctx, req)
	if err != nil {
		logError("Failed to permanently delete role via command client", err, zap.Int("role_id", roleID), zap.Error(err))
		return errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseRoleDelete(res)
	logSuccess("Successfully permanently deleted role", zap.String("handler", "DeletePermanent"), zap.Int("role_id", roleID))

	return c.JSON(http.StatusOK, so)
}

// RestoreAll godoc.
// @Summary Restore all soft-deleted roles
// @Tags Role
// @Security Bearer
// @Description Restore all soft-deleted roles.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Restored roles data"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all roles"
// @Router /api/role/restore/all [post]
func (h *roleHandleApi) RestoreAll(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "RestoreAll"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to restore all trashed roles", zap.String("handler", "RestoreAll"))

	res, err := h.commandClient.RestoreAllRole(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to restore all roles via command client", err, zap.Error(err))
		return errors.ErrApiFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseRoleAll(res)
	logSuccess("Successfully restored all trashed roles", zap.String("handler", "RestoreAll"))

	return c.JSON(http.StatusOK, so)
}

// DeleteAllPermanent godoc.
// @Summary Permanently delete all roles
// @Tags Role
// @Security Bearer
// @Description Permanently delete all roles.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Permanently deleted roles data"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all roles permanently"
// @Router /api/role/permanent/all [delete]
func (h *roleHandleApi) DeleteAllPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "role"); err != nil {
		return err
	}

	const method = "DeleteAllPermanent"
	ctx := c.Request().Context()
	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to permanently delete all trashed roles", zap.String("handler", "DeleteAllPermanent"))

	res, err := h.commandClient.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to permanently delete all roles via command client", err, zap.Error(err))
		return errors.ErrApiFailedDeleteAll(c)
	}

	so := h.mapping.ToApiResponseRoleAll(res)
	logSuccess("Successfully permanently deleted all trashed roles", zap.String("handler", "DeleteAllPermanent"))

	return c.JSON(http.StatusOK, so)
}
