package grpcclient

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpclientmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/response/grpclient"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/middlewares"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	pbcommon "github.com/MamangRust/simple_microservice_ecommerce_pb/common"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type UserGrpcClientHandler interface {
	FindByEmail(ctx context.Context, email string) (*response.ApiResponseUser, *response.ErrorResponse)
	FindById(ctx context.Context, id int32) (*response.ApiResponseUser, *response.ErrorResponse)
	FindByEmailAndVerify(ctx context.Context, email string) (*response.ApiResponseUserWithPassword, *response.ErrorResponse)
	CreateUser(ctx context.Context, req *requests.CreateUserRequest) (*response.ApiResponseUser, *response.ErrorResponse)
	UpdateUserIsVerified(ctx context.Context, req *requests.UpdateUserVerifiedRequest) (*response.ApiResponseUser, *response.ErrorResponse)
	UpdateUserPassword(ctx context.Context, req *requests.UpdateUserPasswordRequest) (*response.ApiResponseUser, *response.ErrorResponse)
	FindVerificationCode(ctx context.Context, code string) (*response.ApiResponseUser, *response.ErrorResponse)
}

type userGrpcClientHandler struct {
	clientQuery   pbuser.UserQueryServiceClient
	clientCommand pbuser.UserCommandServiceClient
	errorHandler  middlewares.GRPCErrorHandling
	logger        logger.LoggerInterface
	mapper        grpclientmapper.UserClientResponseMapper
}

func NewUserGrpcClientHandler(clientQuery pbuser.UserQueryServiceClient,
	clientCommand pbuser.UserCommandServiceClient, logger logger.LoggerInterface, errorHandler middlewares.GRPCErrorHandling) UserGrpcClientHandler {
	mapper := grpclientmapper.NewUserClientResponseMapper()

	return &userGrpcClientHandler{
		clientQuery:   clientQuery,
		clientCommand: clientCommand,
		errorHandler:  errorHandler,
		logger:        logger,
		mapper:        mapper,
	}
}

func (h *userGrpcClientHandler) FindByEmail(ctx context.Context, email string) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuser.FindByEmailUserRequest{Email: email}

	res, err := h.clientQuery.FindByEmail(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}

		h.logger.Error("gRPC client error when finding user by email", zap.Error(err))
		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return h.mapper.ToApiResponseUser(res), nil
}

func (h *userGrpcClientHandler) FindById(ctx context.Context, id int32) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuser.FindByIdUserRequest{Id: id}

	res, err := h.clientQuery.FindById(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}

		h.logger.Error("gRPC client error when finding user by email", zap.Int32("id", id), zap.Error(err))

		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return h.mapper.ToApiResponseUser(res), nil
}

func (h *userGrpcClientHandler) FindByEmailAndVerify(ctx context.Context, email string) (*response.ApiResponseUserWithPassword, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuser.FindByEmailUserRequest{Email: email}

	res, err := h.clientQuery.FindByEmailAndVerify(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}

		h.logger.Error("gRPC client error when finding user for verification", zap.String("email", email), zap.Error(err))
		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	h.logger.Info("FindByEmailAndVerify response Data",
		zap.String("email", email),
		zap.Any("data_json", res.Data),
	)

	return h.mapper.ToApiResponseUserWithPassword(res), nil
}

func (h *userGrpcClientHandler) CreateUser(ctx context.Context, req *requests.CreateUserRequest) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	reqPb := &pbcommon.CreateUserRequest{
		Firstname:       req.FirstName,
		Lastname:        req.LastName,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		IsVerified:      req.IsVerified,
		VerifiedCode:    req.VerifiedCode,
	}

	res, err := h.clientCommand.CreateUser(ctx, reqPb)
	if err != nil {
		h.logger.Error("grpc client error when create user", zap.String("email", req.Email), zap.Error(err))

		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return h.mapper.ToApiResponseUser(res), nil
}

func (h *userGrpcClientHandler) UpdateUserIsVerified(ctx context.Context, req *requests.UpdateUserVerifiedRequest) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	reqPb := &pbuser.UpdateUserVerifiedRequest{
		UserId:     int32(req.UserID),
		IsVerified: req.IsVerified,
	}

	res, err := h.clientCommand.UpdateUserIsVerified(ctx, reqPb)
	if err != nil {
		h.logger.Error("grpc client error when update user verified", zap.Int32("userId", int32(req.UserID)), zap.Error(err))

		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return h.mapper.ToApiResponseUser(res), nil
}

func (h *userGrpcClientHandler) UpdateUserPassword(ctx context.Context, req *requests.UpdateUserPasswordRequest) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	if h.clientCommand == nil {
		h.logger.Error("user command client is nil in UpdateUserPassword")
		return nil, response.NewApiErrorResponse("error", "user service unavailable", 503)
	}

	reqPb := &pbuser.UpdateUserPasswordRequest{
		UserId:   int32(req.UserID),
		Password: req.Password,
	}

	res, err := h.clientCommand.UpdateUserPassword(ctx, reqPb)
	if err != nil {
		h.logger.Error("grpc client error when update user password", zap.Int32("userId", int32(req.UserID)), zap.Error(err))

		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return h.mapper.ToApiResponseUser(res), nil
}

func (h *userGrpcClientHandler) FindVerificationCode(ctx context.Context, code string) (*response.ApiResponseUser, *response.ErrorResponse) {
	if errResp := h.errorHandler.ValidateClient(h.clientQuery, "user-service"); errResp != nil {
		return nil, errResp
	}

	req := &pbuser.VerifyCodeRequest{Code: code}

	res, err := h.clientQuery.FindVerificationCode(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}

		h.logger.Error("gRPC client error when finding verify code ", zap.Error(err))

		return nil, h.errorHandler.HandleGRPCError(err, "user-service")
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	return h.mapper.ToApiResponseUser(res), nil
}
