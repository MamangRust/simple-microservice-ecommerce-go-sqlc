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
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleApi struct {
	queryClient    pb.UserQueryServiceClient
	commandClient  pb.UserCommandServiceClient
	logger         logger.LoggerInterface
	grpcmiddleware middlewares.GRPCErrorHandlingMiddleware
	mapping        mapper.UserResponseMapper
	mencache       mencache.UserCache
	observability  observability.TraceLoggerObservability
}

func NewUserHandleApi(router *echo.Echo, queryClient pb.UserQueryServiceClient, commandClient pb.UserCommandServiceClient, logger logger.LoggerInterface, grpcmiddleware middlewares.GRPCErrorHandlingMiddleware, mencache mencache.UserCache) *userHandleApi {
	observability, _ := observability.NewObservability("user-service", logger)

	userHandler := &userHandleApi{
		queryClient:    queryClient,
		commandClient:  commandClient,
		logger:         logger,
		grpcmiddleware: grpcmiddleware,
		mapping:        mapper.NewUserResponseMapper(),
		mencache:       mencache,
		observability:  observability,
	}

	routerUser := router.Group("/api/user")

	routerUser.GET("", userHandler.FindAllUser)
	routerUser.GET("/:id", userHandler.FindById)
	routerUser.GET("/active", userHandler.FindByActive)
	routerUser.GET("/trashed", userHandler.FindByTrashed)

	routerUser.POST("/update/:id", userHandler.Update)

	routerUser.POST("/trashed/:id", userHandler.TrashedUser)
	routerUser.POST("/restore/:id", userHandler.RestoreUser)
	routerUser.DELETE("/permanent/:id", userHandler.DeleteUserPermanent)

	routerUser.POST("/restore/all", userHandler.RestoreAllUser)
	routerUser.POST("/permanent/all", userHandler.DeleteAllUserPermanent)

	return userHandler
}

// @Security Bearer
// @Summary Find all users
// @Tags User
// @Description Retrieve a list of all users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUser "List of users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user [get]
func (h *userHandleApi) FindAllUser(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "user"); err != nil {
		return err
	}
	const method = "FindAllUser"
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

	logSuccess("Attempting to find all users",
		zap.String("handler", "FindAllUser"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedUsers(ctx, reqCache); found {
		logSuccess("All users found in cache", zap.String("handler", "FindAllUser"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindAll(ctx, req)
	if err != nil {
		logError("Failed to find all users via query client", err, zap.Error(err))
		return errors.ErrApiUserFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationUser(res)
	h.mencache.SetCachedUsers(ctx, reqCache, so)

	logSuccess("Successfully found all users", zap.String("handler", "FindAllUser"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find user by ID
// @Tags User
// @Description Retrieve a user by ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUser "User data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user/{id} [get]
func (h *userHandleApi) FindById(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "user"); err != nil {
		return err
	}
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid user ID provided in FindById", err, zap.String("handler", "FindById"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiUserInvalidId(c)
	}

	logSuccess("Attempting to find user by ID", zap.String("handler", "FindById"), zap.Int("user_id", id))

	if data, found := h.mencache.GetCachedUserById(ctx, id); found {
		logSuccess("User found in cache", zap.String("handler", "FindById"), zap.Int("user_id", id))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindByIdUserRequest{Id: int32(id)}
	user, err := h.queryClient.FindById(ctx, req)
	if err != nil {
		logError("Failed to find user by ID via query client", err, zap.Int("user_id", id), zap.Error(err))
		return errors.ErrApiUserNotFound(c)
	}

	so := h.mapping.ToApiResponseUser(user)
	h.mencache.SetCachedUserById(ctx, so)

	logSuccess("Successfully found user by ID", zap.String("handler", "FindById"), zap.Int("user_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active users
// @Tags User
// @Description Retrieve a list of active users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUserDeleteAt "List of active users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user/active [get]
func (h *userHandleApi) FindByActive(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "user"); err != nil {
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

	logSuccess("Attempting to find active users",
		zap.String("handler", "FindByActive"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedUserActive(ctx, reqCache); found {
		logSuccess("Active users found in cache", zap.String("handler", "FindByActive"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByActive(ctx, req)
	if err != nil {
		logError("Failed to find active users via query client", err, zap.Error(err))
		return errors.ErrApiUserFailedFindActive(c)
	}

	so := h.mapping.ToApiResponsePaginationUserDeleteAt(res)
	h.mencache.SetCachedUserActive(ctx, reqCache, so)

	logSuccess("Successfully found active users", zap.String("handler", "FindByActive"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed user records.
// @Summary Retrieve trashed users
// @Tags User
// @Description Retrieve a list of trashed user records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUserDeleteAt "List of trashed user data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user/trashed [get]
func (h *userHandleApi) FindByTrashed(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.queryClient, "user"); err != nil {
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

	logSuccess("Attempting to find trashed users",
		zap.String("handler", "FindByTrashed"),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqCache := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	if data, found := h.mencache.GetCachedUserTrashed(ctx, reqCache); found {
		logSuccess("Trashed users found in cache", zap.String("handler", "FindByTrashed"))
		return c.JSON(http.StatusOK, data)
	}

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.queryClient.FindByTrashed(ctx, req)
	if err != nil {
		logError("Failed to find trashed users via query client", err, zap.Error(err))
		return errors.ErrApiUserFailedFindTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationUserDeleteAt(res)
	h.mencache.SetCachedUserTrashed(ctx, reqCache, so)

	logSuccess("Successfully found trashed users", zap.String("handler", "FindByTrashed"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing user record.
// @Summary Update an existing user
// @Tags User
// @Description Update an existing user record with the provided details
// @Accept json
// @Produce json
// @Param UpdateUserRequest body requests.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.ApiResponseUser "Successfully updated user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update user"
// @Router /api/user/update/{id} [post]
func (h *userHandleApi) Update(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "user"); err != nil {
		return err
	}
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid user ID provided in Update", err, zap.String("handler", "Update"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiUserInvalidId(c)
	}

	var body requests.UpdateUserRequest
	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateUser request", err, zap.Error(err))
		return errors.ErrApiUserBindUpdateUser(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate UpdateUser request", err, zap.Int("user_id", idInt), zap.String("email", body.Email), zap.Error(err))
		return errors.ErrApiUserValidateUpdateUser(c)
	}

	logSuccess("Attempting to update user", zap.String("handler", "Update"), zap.Int("user_id", idInt), zap.String("email", body.Email))

	req := &pb.UpdateUserRequest{
		Id:              int32(idInt),
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.commandClient.UpdateUser(ctx, req)
	if err != nil {
		logError("Failed to update user via command client", err, zap.Int("user_id", idInt), zap.Error(err))
		return errors.ErrApiUserFailedUpdateUser(c)
	}

	so := h.mapping.ToApiResponseUser(res)
	logSuccess("Successfully updated user", zap.String("handler", "Update"), zap.Int("user_id", idInt))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedUser retrieves a trashed user record by its ID.
// @Summary Retrieve a trashed user
// @Tags User
// @Description Retrieve a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDeleteAt "Successfully retrieved trashed user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed user"
// @Router /api/user/trashed/{id} [get]
func (h *userHandleApi) TrashedUser(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "user"); err != nil {
		return err
	}
	const method = "TrashedUser"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid user ID provided in TrashedUser", err, zap.String("handler", "TrashedUser"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiUserInvalidId(c)
	}

	logSuccess("Attempting to trash user", zap.String("handler", "TrashedUser"), zap.Int("user_id", id))

	req := &pb.FindByIdUserRequest{Id: int32(id)}
	user, err := h.commandClient.TrashedUser(ctx, req)
	if err != nil {
		logError("Failed to trash user via command client", err, zap.Int("user_id", id), zap.Error(err))
		return errors.ErrApiUserFailedTrashedUser(c)
	}

	so := h.mapping.ToApiResponseUserDeleteAt(user)
	logSuccess("Successfully trashed user", zap.String("handler", "TrashedUser"), zap.Int("user_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreUser restores a user record from the trash by its ID.
// @Summary Restore a trashed user
// @Tags User
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDeleteAt "Successfully restored user"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user/restore/{id} [post]
func (h *userHandleApi) RestoreUser(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "user"); err != nil {
		return err
	}

	const method = "RestoreUser"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid user ID provided in RestoreUser", err, zap.String("handler", "RestoreUser"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiUserInvalidId(c)
	}

	logSuccess("Attempting to restore user", zap.String("handler", "RestoreUser"), zap.Int("user_id", id))

	req := &pb.FindByIdUserRequest{Id: int32(id)}

	user, err := h.commandClient.RestoreUser(ctx, req)

	if err != nil {
		logError("Failed to restore user via command client", err, zap.Int("user_id", id), zap.Error(err))
		return errors.ErrApiUserFailedRestoreUser(c)
	}

	so := h.mapping.ToApiResponseUserDeleteAt(user)
	logSuccess("Successfully restored user", zap.String("handler", "RestoreUser"), zap.Int("user_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteUserPermanent permanently deletes a user record by its ID.
// @Summary Permanently delete a user
// @Tags User
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user:"
// @Router /api/user/delete/{id} [delete]
func (h *userHandleApi) DeleteUserPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "user"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	const method = "DeleteUserPermanent"

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid user ID provided in DeleteUserPermanent", err, zap.String("handler", "DeleteUserPermanent"), zap.String("param_id", c.Param("id")))
		return errors.ErrApiUserInvalidId(c)
	}

	logSuccess("Attempting to permanently delete user", zap.String("handler", "DeleteUserPermanent"), zap.Int("user_id", id))

	req := &pb.FindByIdUserRequest{Id: int32(id)}
	user, err := h.commandClient.DeleteUserPermanent(ctx, req)
	if err != nil {
		logError("Failed to permanently delete user via command client", err, zap.Int("user_id", id), zap.Error(err))
		return errors.ErrApiUserFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseUserDelete(user)
	logSuccess("Successfully permanently deleted user", zap.String("handler", "DeleteUserPermanent"), zap.Int("user_id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreUser restores a user record from the trash by its ID.
// @Summary Restore a trashed user
// @Tags User
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserAll "Successfully restored user all"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user/restore/all [post]
func (h *userHandleApi) RestoreAllUser(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "user"); err != nil {
		return err
	}
	const method = "RestoreAllUser"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to restore all trashed users", zap.String("handler", "RestoreAllUser"))

	res, err := h.commandClient.RestoreAllUser(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to restore all users via command client", err, zap.Error(err))
		return errors.ErrApiUserFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseUserAll(res)
	logSuccess("Successfully restored all trashed users", zap.String("handler", "RestoreAllUser"))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteUserPermanent permanently deletes a user record by its ID.
// @Summary Permanently delete a user
// @Tags User
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user:"
// @Router /api/user/delete/all [post]
func (h *userHandleApi) DeleteAllUserPermanent(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.commandClient, "user"); err != nil {
		return err
	}

	const method = "DeleteAllUserPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.observability.StartTracingAndLogging(ctx, method)

	defer func() { end() }()

	logSuccess("Attempting to permanently delete all trashed users", zap.String("handler", "DeleteAllUserPermanent"))

	res, err := h.commandClient.DeleteAllUserPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to permanently delete all users via command client", err, zap.Error(err))
		return errors.ErrApiUserFailedDeleteAll(c)
	}

	so := h.mapping.ToApiResponseUserAll(res)
	logSuccess("Successfully permanently deleted all trashed users", zap.String("handler", "DeleteAllUserPermanent"))

	return c.JSON(http.StatusOK, so)
}
