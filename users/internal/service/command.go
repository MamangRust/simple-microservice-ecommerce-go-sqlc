package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/user/internal/grpc_client"
	responsemapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userCommandDeps struct {
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleGrpcClient        grpcclient.RoleGrpcClientHandler
	userRoleGrpcClient    grpcclient.UserRoleGrpcClientHandler
	logger                logger.LoggerInterface
	hash                  hash.HashPassword
	mapper                responsemapper.UserCommandResponseMapper
}

type userCommandService struct {
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleGrpcClient        grpcclient.RoleGrpcClientHandler
	userRoleGrpcClient    grpcclient.UserRoleGrpcClientHandler
	logger                logger.LoggerInterface
	hash                  hash.HashPassword
	mapper                responsemapper.UserCommandResponseMapper
}

func NewUserCommandService(
	params *userCommandDeps,
) UserCommandService {
	return &userCommandService{
		userQueryRepository:   params.userQueryRepository,
		userCommandRepository: params.userCommandRepository,
		roleGrpcClient:        params.roleGrpcClient,
		hash:                  params.hash,
		userRoleGrpcClient:    params.userRoleGrpcClient,
		logger:                params.logger,
		mapper:                params.mapper,
	}
}

func (s *userCommandService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("Attempting to create user", zap.String("email", request.Email))

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingUser = nil
		} else {
			s.logger.Error("Failed to check for existing user", zap.Error(err), zap.String("email", request.Email))
			return nil, response.NewErrorResponse("failed to check for existing user", 500)
		}
	}
	if existingUser != nil {
		s.logger.Info("User creation failed: email already exists", zap.String("email", request.Email))
		return nil, response.NewErrorResponse("user with this email already exists", 409)
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		s.logger.Error("Failed to hash password for new user", zap.Error(err), zap.String("email", request.Email))
		return nil, response.NewErrorResponse("failed to process password", 500)
	}
	request.Password = passwordHash

	const defaultRoleName = "admin"
	s.logger.Info("Attempting to find default role for new user", zap.String("role_name", defaultRoleName))
	role, errResp := s.roleGrpcClient.FindByName(ctx, defaultRoleName)
	if errResp != nil {
		st, ok := status.FromError(errResp.ToGRPCError())
		if ok && st.Code() == codes.NotFound {
			s.logger.Error("Default role not found during user creation", zap.String("role_name", defaultRoleName))
			return nil, response.NewErrorResponse("default role 'ADMIN' not found", 404)
		}
		s.logger.Error("Failed to communicate with role service", zap.String("error_message", errResp.Message))
		return nil, response.NewErrorResponse("failed to communicate with role service", 500)
	}

	newUser, err := s.userCommandRepository.CreateUser(ctx, request)
	if err != nil {
		s.logger.Error("Failed to create user in database", zap.Error(err), zap.String("email", request.Email))
		return nil, response.NewErrorResponse("failed to create user", 500)
	}

	s.logger.Info("Attempting to assign role to new user", zap.Int("user_id", newUser.ID), zap.Int("role_id", role.Data.ID))
	_, errResp = s.userRoleGrpcClient.AssignUserRole(ctx, &requests.CreateUserRoleRequest{
		UserId: newUser.ID,
		RoleId: role.Data.ID,
	})
	if errResp != nil {
		s.logger.Error("Failed to assign role to new user", zap.String("error_message", errResp.Message), zap.Int("user_id", newUser.ID))
		return nil, response.NewErrorResponse("failed to assign role to user", 500)
	}

	s.logger.Info("User created successfully", zap.Int("user_id", newUser.ID), zap.String("email", newUser.Email))
	so := s.mapper.ToUserResponse(newUser)
	return so, nil
}

func (s *userCommandService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("Attempting to update user", zap.Int("user_id", *request.UserID))

	existingUser, err := s.userQueryRepository.FindById(ctx, *request.UserID)
	if err != nil {
		s.logger.Error("Failed to find user for update", zap.Error(err), zap.Int("user_id", *request.UserID))
		return nil, response.NewErrorResponse("user not found", 404)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		s.logger.Info("Checking for duplicate email during user update", zap.String("new_email", request.Email))
		duplicateUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
		if err != nil {
			s.logger.Error("Failed to check for duplicate email", zap.Error(err), zap.String("new_email", request.Email))
			return nil, response.NewErrorResponse("failed to check for duplicate email", 500)
		}
		if duplicateUser != nil {
			s.logger.Info("User update failed: new email is already in use", zap.String("new_email", request.Email))
			return nil, response.NewErrorResponse("email is already in use by another user", 409)
		}
	}

	if request.Password != "" {
		s.logger.Info("Hashing new password for user update", zap.Int("user_id", *request.UserID))
		hash, err := s.hash.HashPassword(request.Password)
		if err != nil {
			s.logger.Error("Failed to hash new password", zap.Error(err), zap.Int("user_id", *request.UserID))
			return nil, response.NewErrorResponse("failed to process new password", 500)
		}
		existingUser.Password = hash
	}

	const defaultRoleName = "admin"
	s.logger.Info("Attempting to find default role for user update", zap.String("role_name", defaultRoleName))
	resrole, errResp := s.roleGrpcClient.FindByName(ctx, defaultRoleName)
	if errResp != nil {
		st, ok := status.FromError(errResp.ToGRPCError())
		if ok && st.Code() == codes.NotFound {
			s.logger.Error("Default role not found during user update", zap.String("role_name", defaultRoleName))
			return nil, response.NewErrorResponse("default role 'ADMIN' not found", 404)
		}
		s.logger.Error("Failed to communicate with role service during user update", zap.String("error_message", errResp.Message))
		return nil, response.NewErrorResponse("failed to communicate with role service", 500)
	}

	s.logger.Info("Attempting to update user role", zap.Int("user_id", *request.UserID), zap.Int("role_id", resrole.Data.ID))
	_, errResp = s.userRoleGrpcClient.UpdateUserRole(ctx, &requests.UpdateUserRoleRequest{
		UserId: *request.UserID,
		RoleId: resrole.Data.ID,
	})
	if errResp != nil {
		s.logger.Error("Failed to update user role", zap.String("error_message", errResp.Message), zap.Int("user_id", *request.UserID))
		return nil, response.NewErrorResponse("failed to update user role", 500)
	}

	res, err := s.userCommandRepository.UpdateUser(ctx, request)
	if err != nil {
		s.logger.Error("Failed to update user data in database", zap.Error(err), zap.Int("user_id", *request.UserID))
		return nil, response.NewErrorResponse("failed to update user data", 500)
	}

	s.logger.Info("User updated successfully", zap.Int("user_id", res.ID))
	so := s.mapper.ToUserResponse(res)
	return so, nil
}

func (s *userCommandService) UpdateUserIsVerified(ctx context.Context, request *requests.UpdateUserVerifiedRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("Attempting to update user verification status",
		zap.Int("user_id", request.UserID),
		zap.Bool("is_verified", request.IsVerfied),
	)

	_, err := s.userQueryRepository.FindById(ctx, request.UserID)
	if err != nil {
		s.logger.Error("Failed to find user for verification update", zap.Error(err), zap.Int("user_id", request.UserID))
		return nil, response.NewErrorResponse("user not found", 404)
	}

	res, err := s.userCommandRepository.UpdateUserIsVerified(ctx, request)
	if err != nil {
		s.logger.Error("Failed to update user verification status in database", zap.Error(err), zap.Int("user_id", request.UserID))
		return nil, response.NewErrorResponse("failed to update user data", 500)
	}

	s.logger.Info("User verification status updated successfully", zap.Int("user_id", res.ID))
	so := s.mapper.ToUserResponse(res)
	return so, nil
}

func (s *userCommandService) UpdateUserPassword(ctx context.Context, request *requests.UpdateUserPasswordRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("Attempting to update user password", zap.Int("user_id", request.UserID))

	_, err := s.userQueryRepository.FindById(ctx, request.UserID)
	if err != nil {
		s.logger.Error("Failed to find user for password update", zap.Error(err), zap.Int("user_id", request.UserID))
		return nil, response.NewErrorResponse("user not found", 404)
	}

	res, err := s.userCommandRepository.UpdateUserPassword(ctx, request)
	if err != nil {
		s.logger.Error("Failed to update user password in database", zap.Error(err), zap.Int("user_id", request.UserID))
		return nil, response.NewErrorResponse("failed to update user data", 500)
	}

	s.logger.Info("User password updated successfully", zap.Int("user_id", res.ID))
	so := s.mapper.ToUserResponse(res)
	return so, nil
}

func (s *userCommandService) TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Info("Attempting to trash user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.TrashedUser(ctx, user_id)
	if err != nil {
		s.logger.Error("Failed to trash user in database", zap.Error(err), zap.Int("user_id", user_id))
		return nil, response.NewErrorResponse("failed to trash user", 500)
	}

	s.logger.Info("User trashed successfully", zap.Int("user_id", res.ID))
	so := s.mapper.ToUserResponseDeleteAt(res)
	return so, nil
}

func (s *userCommandService) RestoreUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Info("Attempting to restore user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.RestoreUser(ctx, user_id)
	if err != nil {
		s.logger.Error("Failed to restore user in database", zap.Error(err), zap.Int("user_id", user_id))
		return nil, response.NewErrorResponse("failed to restore user", 500)
	}

	s.logger.Info("User restored successfully", zap.Int("user_id", res.ID))
	so := s.mapper.ToUserResponseDeleteAt(res)
	return so, nil
}

func (s *userCommandService) DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse) {
	s.logger.Info("Attempting to permanently delete user", zap.Int("user_id", user_id))

	_, err := s.userCommandRepository.DeleteUserPermanent(ctx, user_id)
	if err != nil {
		s.logger.Error("Failed to permanently delete user from database", zap.Error(err), zap.Int("user_id", user_id))
		return false, response.NewErrorResponse("failed to delete user permanently", 500)
	}

	s.logger.Info("User permanently deleted successfully", zap.Int("user_id", user_id))
	return true, nil
}

func (s *userCommandService) RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse) {
	s.logger.Info("Attempting to restore all trashed users")

	_, err := s.userCommandRepository.RestoreAllUser(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all users from database", zap.Error(err))
		return false, response.NewErrorResponse("failed to restore all users", 500)
	}

	s.logger.Info("All trashed users restored successfully")
	return true, nil
}

func (s *userCommandService) DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	s.logger.Info("Attempting to permanently delete all trashed users")

	_, err := s.userCommandRepository.DeleteAllUserPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all users from database", zap.Error(err))
		return false, response.NewErrorResponse("failed to delete all users permanently", 500)
	}

	s.logger.Info("All trashed users permanently deleted successfully")
	return true, nil
}
