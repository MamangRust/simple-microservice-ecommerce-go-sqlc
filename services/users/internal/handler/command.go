package handler

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	protomapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/proto"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/service"
	usergrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors/user_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	pbcommon "github.com/MamangRust/simple_microservice_ecommerce_pb/common"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userCommandHandleGrpc struct {
	pbuser.UnimplementedUserCommandServiceServer
	userCommandService service.UserCommandService
	logger             logger.LoggerInterface
	mapper             protomapper.UserCommandProtoMapper
}

func NewUserCommandHandleGrpc(command service.UserCommandService, logger logger.LoggerInterface) pbuser.UserCommandServiceServer {
	return &userCommandHandleGrpc{
		userCommandService: command,
		logger:             logger,
		mapper:             protomapper.NewUserCommandProtoMapper(),
	}
}

func (s *userCommandHandleGrpc) CreateUser(ctx context.Context, request *pbcommon.CreateUserRequest) (*pbuser.ApiResponseUser, error) {
	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
		VerifiedCode:    request.VerifiedCode,
		IsVerified:      request.IsVerified,
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("Invalid request data for CreateUser", zap.Error(err))
		return nil, usergrpcerrors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.CreateUser(ctx, req)

	if err != nil {
		s.logger.Error("Failed to create user in service",
			zap.String("error_message", err.Message),
			zap.String("email", req.Email),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User successfully created",
		zap.Int("userId", user.ID),
		zap.String("email", user.Email),
	)

	so := s.mapper.ToProtoResponseUser("success", "Successfully created user", user)
	return so, nil
}

func (s *userCommandHandleGrpc) UpdateUser(ctx context.Context, request *pbuser.UpdateUserRequest) (*pbuser.ApiResponseUser, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid UserID received on update request", zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	req := &requests.UpdateUserRequest{
		UserID:          &id,
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("Invalid request data for UpdateUser", zap.Error(err), zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.UpdateUser(ctx, req)

	if err != nil {
		s.logger.Error("Failed to update user in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User successfully updated", zap.Int("userId", id))

	so := s.mapper.ToProtoResponseUser("success", "Successfully updated user", user)
	return so, nil
}

func (s *userCommandHandleGrpc) UpdateUserIsVerified(ctx context.Context, request *pbuser.UpdateUserVerifiedRequest) (*pbuser.ApiResponseUser, error) {
	id := int(request.GetUserId())

	if id == 0 {
		s.logger.Warn("Invalid UserID received on update verification request", zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	res, err := s.userCommandService.UpdateUserIsVerified(ctx, &requests.UpdateUserVerifiedRequest{
		UserID:    int(request.GetUserId()),
		IsVerfied: request.IsVerified,
	})

	if err != nil {
		s.logger.Error("Failed to update user verification status in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User verification status successfully updated",
		zap.Int("userId", id),
		zap.Bool("isVerified", request.IsVerified),
	)

	so := s.mapper.ToProtoResponseUser("success", "Successfully updated user", res)
	return so, nil
}

func (s *userCommandHandleGrpc) UpdateUserPassword(ctx context.Context, request *pbuser.UpdateUserPasswordRequest) (*pbuser.ApiResponseUser, error) {
	id := int(request.GetUserId())

	if id == 0 {
		s.logger.Warn("Invalid UserID received on update password request", zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	res, err := s.userCommandService.UpdateUserPassword(ctx, &requests.UpdateUserPasswordRequest{
		UserID:   int(request.GetUserId()),
		Password: request.GetPassword(),
	})

	if err != nil {
		s.logger.Error("Failed to update user password in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User password successfully updated", zap.Int("userId", id))

	so := s.mapper.ToProtoResponseUser("success", "Successfully updated user", res)
	return so, nil
}

func (s *userCommandHandleGrpc) TrashedUser(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid UserID received on trash request", zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.TrashedUser(ctx, id)

	if err != nil {
		s.logger.Error("Failed to trash user in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User successfully trashed", zap.Int("userId", id))

	so := s.mapper.ToProtoResponseUserDeleteAt("success", "Successfully trashed user", user)
	return so, nil
}

func (s *userCommandHandleGrpc) RestoreUser(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid UserID received on restore request", zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.RestoreUser(ctx, id)

	if err != nil {
		s.logger.Error("Failed to restore user in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User successfully restored", zap.Int("userId", id))

	so := s.mapper.ToProtoResponseUserDeleteAt("success", "Successfully restored user", user)
	return so, nil
}

func (s *userCommandHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUserDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid UserID received on permanent delete request", zap.Int("userId", id))
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	_, err := s.userCommandService.DeleteUserPermanent(ctx, id)

	if err != nil {
		s.logger.Error("Failed to permanently delete user in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("User successfully permanently deleted", zap.Int("userId", id))

	so := s.mapper.ToProtoResponseUserDelete("success", "Successfully deleted user permanently")
	return so, nil
}

func (s *userCommandHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pbuser.ApiResponseUserAll, error) {
	_, err := s.userCommandService.RestoreAllUser(ctx)

	if err != nil {
		s.logger.Error("Failed to restore all users in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("All trashed users successfully restored")

	so := s.mapper.ToProtoResponseUserAll("success", "Successfully restored all users")
	return so, nil
}

func (s *userCommandHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pbuser.ApiResponseUserAll, error) {
	_, err := s.userCommandService.DeleteAllUserPermanent(ctx)

	if err != nil {
		s.logger.Error("Failed to permanently delete all users in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("All trashed users successfully permanently deleted")

	so := s.mapper.ToProtoResponseUserAll("success", "Successfully deleted all user permanently")
	return so, nil
}
