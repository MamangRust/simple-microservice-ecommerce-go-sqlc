package handler

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/mapper"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/errors"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
	pbcommon "github.com/MamangRust/simple_microservice_ecommerce_pb/common"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type authHandleApi struct {
	client         pb.AuthServiceClient
	logger         logger.LoggerInterface
	grpcmiddleware middlewares.GRPCErrorHandlingMiddleware
	mapping        mapper.AuthResponseMapper
}

func NewAuthHandle(router *echo.Echo, client pb.AuthServiceClient, logger logger.LoggerInterface, grpcmiddleware middlewares.GRPCErrorHandlingMiddleware) *authHandleApi {
	authHandler := &authHandleApi{
		client:         client,
		logger:         logger,
		grpcmiddleware: grpcmiddleware,
		mapping:        mapper.NewAuthResponseMapper(),
	}

	routerAuth := router.Group("/api/auth")

	routerAuth.GET("/verify-code", authHandler.VerifyCode)
	routerAuth.POST("/forgot-password", authHandler.ForgotPassword)
	routerAuth.POST("/reset-password", authHandler.ResetPassword)
	routerAuth.GET("/hello", authHandler.HandleHello)
	routerAuth.POST("/register", authHandler.Register)
	routerAuth.POST("/login", authHandler.Login)
	routerAuth.POST("/refresh-token", authHandler.RefreshToken)
	routerAuth.GET("/me", authHandler.GetMe)

	return authHandler
}

// HandleHello godoc
// @Summary Returns a "Hello" message
// @Tags Auth
// @Description Returns a simple "Hello" message for testing purposes.
// @Produce json
// @Success 200 {string} string "Hello"
// @Router /auth/hello [get]
func (h *authHandleApi) HandleHello(c echo.Context) error {
	return c.String(200, "Hello")
}

// VerifyCode godoc
// @Summary Verifies the user using a verification code
// @Tags Auth
// @Description Verifies the user's email using the verification code provided in the query parameter.
// @Produce json
// @Param verify_code query string true "Verification Code"
// @Success 200 {object} response.ApiResponseVerifyCode
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/verify-code [get]
func (h *authHandleApi) VerifyCode(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	verifyCode := c.QueryParam("verify_code")
	h.logger.Info("Attempting to verify code", zap.String("handler", "VerifyCode"), zap.String("verify_code", verifyCode))

	ctx := c.Request().Context()

	res, err := h.client.VerifyCode(ctx, &pb.VerifyCodeRequest{
		Code: verifyCode,
	})

	if err != nil {
		h.logger.Error("Failed to verify code via client", zap.Error(err))
		return errors.ErrApiVerifyCode(c)
	}

	resp := h.mapping.ToResponseVerifyCode(res)
	h.logger.Info("Successfully verified code", zap.String("handler", "VerifyCode"))

	return c.JSON(http.StatusOK, resp)
}

// ForgotPassword godoc
// @Summary Sends a reset token to the user's email
// @Tags Auth
// @Description Initiates password reset by sending a reset token to the provided email.
// @Accept json
// @Produce json
// @Param request body requests.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} response.ApiResponseForgotPassword
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *authHandleApi) ForgotPassword(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	var body requests.ForgotPasswordRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Warn("Failed to bind ForgotPassword request", zap.Error(err))
		return errors.ErrBindForgotPassword(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Warn("Failed to validate ForgotPassword request", zap.String("email", body.Email), zap.Error(err))
		return errors.ErrValidateForgotPassword(c)
	}

	h.logger.Info("Attempting to process forgot password request", zap.String("handler", "ForgotPassword"), zap.String("email", body.Email))

	res, err := h.client.ForgotPassword(ctx, &pb.ForgotPasswordRequest{
		Email: body.Email,
	})
	if err != nil {
		h.logger.Error("Failed to process forgot password via client", zap.String("email", body.Email), zap.Error(err))
		return errors.ErrApiForgotPassword(c)
	}

	resp := h.mapping.ToResponseForgotPassword(res)
	h.logger.Info("Successfully processed forgot password request", zap.String("handler", "ForgotPassword"), zap.String("email", body.Email))

	return c.JSON(http.StatusOK, resp)
}

// ResetPassword godoc
// @Summary Resets the user's password using a reset token
// @Tags Auth
// @Description Allows user to reset their password using a valid reset token.
// @Accept json
// @Produce json
// @Param request body requests.CreateResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} response.ApiResponseResetPassword
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/reset-password [post]
func (h *authHandleApi) ResetPassword(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	var body requests.CreateResetPasswordRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Warn("Failed to bind ResetPassword request", zap.Error(err))
		return errors.ErrBindResetPassword(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Warn("Failed to validate ResetPassword request", zap.Error(err))
		return errors.ErrValidateResetPassword(c)
	}

	h.logger.Info("Attempting to reset password", zap.String("handler", "ResetPassword"))

	res, err := h.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		ResetToken:      body.ResetToken,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	})

	if err != nil {
		h.logger.Error("Failed to reset password via client", zap.Error(err))
		return errors.ErrApiResetPassword(c)
	}

	so := h.mapping.ToResponseResetPassword(res)
	h.logger.Info("Successfully reset password", zap.String("handler", "ResetPassword"))

	return c.JSON(http.StatusOK, so)
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Description Registers a new user with the provided details.
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "User registration data"
// @Success 200 {object} response.ApiResponseRegister "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/register [post]
func (h *authHandleApi) Register(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Warn("Failed to bind Register request", zap.Error(err))
		return errors.ErrBindRegister(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Warn("Failed to validate Register request", zap.String("email", body.Email), zap.Error(err))
		return errors.ErrValidateRegister(c)
	}

	h.logger.Info("Attempting to register new user", zap.String("handler", "Register"), zap.String("email", body.Email))

	data := &pbcommon.RegisterRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.RegisterUser(ctx, data)

	if err != nil {
		h.logger.Error("Failed to register user via client", zap.String("email", body.Email), zap.Error(err))
		return errors.ErrApiRegister(c)
	}

	so := h.mapping.ToResponseRegister(res)
	h.logger.Info("Successfully registered new user", zap.String("handler", "Register"), zap.String("email", body.Email))

	return c.JSON(http.StatusOK, so)
}

// Login godoc
// @Summary Authenticate a user
// @Tags Auth
// @Description Authenticates a user using the provided email and password.
// @Accept json
// @Produce json
// @Param request body requests.AuthRequest true "User login credentials"
// @Success 200 {object} response.ApiResponseLogin "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/login [post]
func (h *authHandleApi) Login(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	var body requests.AuthRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Warn("Failed to bind Login request", zap.Error(err))
		return errors.ErrBindLogin(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Warn("Failed to validate Login request", zap.String("email", body.Email), zap.Error(err))
		return errors.ErrValidateLogin(c)
	}

	h.logger.Info("Attempting user login", zap.String("handler", "Login"), zap.String("email", body.Email))

	data := &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	}

	res, err := h.client.LoginUser(ctx, data)

	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			h.logger.Error("Failed to login user via client", zap.String("email", body.Email), zap.String("grpc_code", grpcErr.Code().String()), zap.Error(grpcErr.Err()))
		} else {
			h.logger.Error("Failed to login user via client", zap.String("email", body.Email), zap.Error(err))
		}
		return errors.ErrApiLogin(c)
	}

	mappedResponse := h.mapping.ToResponseLogin(res)
	h.logger.Info("User login successful", zap.String("handler", "Login"), zap.String("email", body.Email))

	return c.JSON(http.StatusOK, mappedResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags Auth
// @Security Bearer
// @Description Refreshes the access token using a valid refresh token.
// @Accept json
// @Produce json
// @Param request body requests.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} response.ApiResponseRefreshToken "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/refresh-token [post]
func (h *authHandleApi) RefreshToken(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	ctx := c.Request().Context()
	var body requests.RefreshTokenRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Warn("Failed to bind RefreshToken request", zap.Error(err))
		return errors.ErrBindRefreshToken(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Warn("Failed to validate RefreshToken request", zap.Error(err))
		return errors.ErrValidateRefreshToken(c)
	}

	h.logger.Info("Attempting to refresh token", zap.String("handler", "RefreshToken"))

	res, err := h.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})

	if err != nil {
		h.logger.Error("Failed to refresh token via client", zap.Error(err))
		return errors.ErrApiRefreshToken(c)
	}

	so := h.mapping.ToResponseRefreshToken(res)
	h.logger.Info("Successfully refreshed token", zap.String("handler", "RefreshToken"))

	return c.JSON(http.StatusOK, so)
}

// GetMe godoc
// @Summary Get current user information
// @Tags Auth
// @Security Bearer
// @Description Retrieves the current user's information using a valid access token from the Authorization header.
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.ApiResponseGetMe "Success"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/me [get]
func (h *authHandleApi) GetMe(c echo.Context) error {
	if err := h.grpcmiddleware.ValidateClient(h.client, "auth"); err != nil {
		return err
	}

	ctx := c.Request().Context()

	userID, ok := c.Get("userId").(int)
	if !ok {
		h.logger.Warn("Invalid user context in GetMe", zap.Any("user_id", c.Get("user_id")))
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user context")
	}

	h.logger.Info("Attempting to get user details", zap.String("handler", "GetMe"), zap.Int("user_id", userID))

	res, err := h.client.GetMe(ctx, &pb.GetMeRequest{
		Id: int32(userID),
	})

	if err != nil {
		h.logger.Error("Failed to get user details via client", zap.Int("user_id", userID), zap.Error(err))
		return errors.ErrApiGetMe(c)
	}

	so := h.mapping.ToResponseGetMe(res)
	h.logger.Info("Successfully retrieved user details", zap.String("handler", "GetMe"), zap.Int("user_id", userID))

	return c.JSON(http.StatusOK, so)
}
