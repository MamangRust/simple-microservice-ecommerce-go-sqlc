package handler

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	protomapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/proto"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
	pbcommon "github.com/MamangRust/simple_microservice_ecommerce_pb/common"
	"go.uber.org/zap"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	service service.Service
	mapping protomapper.AuthProtoMapper
	logger  logger.LoggerInterface
}

func NewAuthHandleGrpc(authService service.Service, logger logger.LoggerInterface) pb.AuthServiceServer {
	return &authHandleGrpc{
		service: authService,
		mapping: protomapper.NewAuthProtoMapper(),
		logger:  logger,
	}
}

func (s *authHandleGrpc) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.ApiResponseVerifyCode, error) {
	_, err := s.service.VerifyCode(ctx, req.Code)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	return s.mapping.ToProtoResponseVerifyCode("success", "Verify code successful"), nil
}

func (s *authHandleGrpc) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ApiResponseForgotPassword, error) {
	_, err := s.service.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	return s.mapping.ToProtoResponseForgotPassword("success", "Forgot password successful"), nil
}

func (s *authHandleGrpc) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseResetPassword, error) {
	_, err := s.service.ResetPassword(ctx, &requests.CreateResetPasswordRequest{
		ResetToken:      req.ResetToken,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return nil, err.ToGRPCError()
	}

	return s.mapping.ToProtoResponseResetPassword("success", "Reset password successful"), nil
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	s.logger.Info("gRPC Login request", zap.String("email", req.Email))

	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.service.Login(ctx, request)
	if err != nil {
		s.logger.Error("Login failed", zap.String("email", req.Email), zap.Any("error", err))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Login successful", zap.String("email", req.Email))

	return s.mapping.ToProtoResponseLogin("success", "Login successful", res), nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	res, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	return s.mapping.ToProtoResponseRefreshToken("success", "Refresh token successful", res), nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	res, err := s.service.GetMe(ctx, int(req.GetId()))
	if err != nil {
		return nil, err.ToGRPCError()
	}

	return s.mapping.ToProtoResponseGetMe("success", "Get user profile successful", res), nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pbcommon.RegisterRequest) (*pb.ApiResponseRegister, error) {
	s.logger.Info("gRPC Register request", zap.String("email", req.Email), zap.String("first_name", req.Firstname), zap.String("last_name", req.Lastname))

	request := &requests.RegisterRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.service.Register(ctx, request)
	if err != nil {
		s.logger.Error("Registration failed", zap.String("email", req.Email), zap.Any("error", err))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Registration successful", zap.String("email", req.Email))

	return s.mapping.ToProtoResponseRegister("success", "Registration successful", res), nil
}
