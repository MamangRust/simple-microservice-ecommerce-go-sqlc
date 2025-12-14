package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/user/internal/grpc_client"
	responsemapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/service"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/user/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userCommandDeps struct {
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleGrpcClient        grpcclient.RoleGrpcClientHandler
	userRoleGrpcClient    grpcclient.UserRoleGrpcClientHandler
	logger                logger.LoggerInterface
	hash                  hash.HashPassword
	mapper                responsemapper.UserCommandResponseMapper
	errorhandler          errorhandler.UserCommandErrorHandler
	cacheQuery            mencache.UserQueryCache
	cacheCommand          mencache.UserCommandCache
}

type userCommandService struct {
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleGrpcClient        grpcclient.RoleGrpcClientHandler
	cacheQuery            mencache.UserQueryCache
	cacheCommand          mencache.UserCommandCache
	userRoleGrpcClient    grpcclient.UserRoleGrpcClientHandler
	logger                logger.LoggerInterface
	hash                  hash.HashPassword
	mapper                responsemapper.UserCommandResponseMapper
	observability         observability.TraceLoggerObservability
	errorhandler          errorhandler.UserCommandErrorHandler
}

func NewUserCommandService(
	params *userCommandDeps,
) UserCommandService {
	observability, _ := observability.NewObservability("user-command-service", params.logger)

	return &userCommandService{
		userQueryRepository:   params.userQueryRepository,
		userCommandRepository: params.userCommandRepository,
		roleGrpcClient:        params.roleGrpcClient,
		hash:                  params.hash,
		userRoleGrpcClient:    params.userRoleGrpcClient,
		logger:                params.logger,
		mapper:                params.mapper,
		cacheQuery:            params.cacheQuery,
		cacheCommand:          params.cacheCommand,
		observability:         observability,
		errorhandler:          params.errorhandler,
	}
}

func (s *userCommandService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "CreateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))
	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingUser = nil
		} else {
			status = "error"

			return s.errorhandler.HandleCreateUserError(err, method, "FAILED_CHECK_EXISTING_USER", span, &status)
		}
	}
	if existingUser != nil {
		status = "error"
		return s.errorhandler.HandleCreateUserError(fmt.Errorf("failed already user exists"), method, "FAILED_USER_ALREADY_EXISTS", span, &status)
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleCreateUserError(err, method, "FAILED_HASH_PASSWORD", span, &status)
	}
	request.Password = passwordHash

	const defaultRoleName = "ROLE_ADMIN"
	s.logger.Info("Attempting to find default role for new user", zap.String("role_name", defaultRoleName))
	role, errResp := s.roleGrpcClient.FindByName(ctx, defaultRoleName)
	if errResp != nil {
		return s.errorhandler.HandleCreateUserError(fmt.Errorf("role not found: %s", errResp.Message), method, "FAILED_FIND_DEFAULT_ROLE", span, &status)
	}

	newUser, err := s.userCommandRepository.CreateUser(ctx, request)

	if err != nil {
		status = "error"

		return s.errorhandler.HandleCreateUserError(err, method, "FAILED_CREATE_USER", span, &status)
	}

	s.logger.Info("Attempting to assign role to new user", zap.Int("user_id", newUser.ID), zap.Int("role_id", role.Data.ID))

	_, errResp = s.userRoleGrpcClient.AssignUserRole(ctx, &requests.CreateUserRoleRequest{
		UserId: newUser.ID,
		RoleId: role.Data.ID,
	})
	if errResp != nil {
		status = "error"

		return s.errorhandler.HandleCreateUserError(fmt.Errorf("failed assign role grpc: %s", errResp.Message), method, "FAILED_ASSIGN_ROLE_GRPC", span, &status)
	}
	so := s.mapper.ToUserResponse(newUser)

	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)

	logSuccess("Successfully created user", zap.Int("user_id", newUser.ID), zap.String("email", newUser.Email))

	return so, nil
}

func (s *userCommandService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "UpdateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", *request.UserID))

	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindById(ctx, *request.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			status = "error"

			notFoundErr := errors.New("user not found for update")
			return s.errorhandler.HandleUpdateUserError(notFoundErr, method, "FAILED_USER_NOT_FOUND", span, &status)
		}
		status = "error"

		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_FIND_USER_FOR_UPDATE", span, &status)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			status = "error"

			return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_CHECK_DUPLICATE_EMAIL", span, &status)
		}
		if duplicateUser != nil {
			status = "error"

			duplicateErr := errors.New("new email is already in use by another user")
			return s.errorhandler.HandleUpdateUserError(duplicateErr, method, "FAILED_EMAIL_ALREADY_EXISTS", span, &status)
		}
	}

	if request.Password != "" {
		passwordHash, err := s.hash.HashPassword(request.Password)
		if err != nil {
			status = "error"

			return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_HASH_PASSWORD_UPDATE", span, &status)
		}
		request.Password = passwordHash
	}

	const defaultRoleName = "ROLE_ADMIN"
	s.logger.Info("Attempting to find default role for new user", zap.String("role_name", defaultRoleName))
	role, errResp := s.roleGrpcClient.FindByName(ctx, defaultRoleName)
	if errResp != nil {
		return s.errorhandler.HandleCreateUserError(fmt.Errorf("role not found: %s", errResp.Message), method, "FAILED_FIND_DEFAULT_ROLE", span, &status)
	}

	_, errResp = s.userRoleGrpcClient.UpdateUserRole(ctx, &requests.UpdateUserRoleRequest{
		UserId: *request.UserID,
		RoleId: role.Data.ID,
	})
	if errResp != nil {
		status = "error"

		return s.errorhandler.HandleUpdateUserError(fmt.Errorf("failed update user role grpc: %s", errResp.Message), method, "FAILED_UPDATE_USER_ROLE_GRPC", span, &status)
	}

	res, err := s.userCommandRepository.UpdateUser(ctx, request)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_UPDATE_USER_REPO", span, &status)
	}
	so := s.mapper.ToUserResponse(res)

	keyId := fmt.Sprintf("user:id:%d", request.UserID)
	keyEmail := fmt.Sprintf("user:email:%s", request.Email)

	s.cacheCommand.DeleteCachedUser(ctx, keyId)
	s.cacheCommand.DeleteCachedUser(ctx, keyEmail)

	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)

	logSuccess("Successfully updated user", zap.Int("user_id", res.ID), zap.String("email", res.Email))

	return so, nil
}

func (s *userCommandService) UpdateUserIsVerified(ctx context.Context, request *requests.UpdateUserVerifiedRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "UpdateUserIsVerified"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", request.UserID), attribute.Bool("is_verified", request.IsVerfied))
	defer func() {
		end(status)
	}()

	_, err := s.userQueryRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_USER_NOT_FOUND", span, &status)
	}

	res, err := s.userCommandRepository.UpdateUserIsVerified(ctx, request)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_UPDATE_VERIFICATION", span, &status)
	}

	userResponse := s.mapper.ToUserResponse(res)

	key := fmt.Sprintf("user:id:%d", request.UserID)

	s.cacheCommand.DeleteCachedUser(ctx, key)
	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)

	logSuccess("Successfully updated user verification status", zap.Int("user_id", res.ID))

	return userResponse, nil
}

func (s *userCommandService) UpdateUserPassword(ctx context.Context, request *requests.UpdateUserPasswordRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "UpdateUserPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", request.UserID))
	defer func() {
		end(status)
	}()

	_, err := s.userQueryRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_USER_NOT_FOUND", span, &status)
	}

	res, err := s.userCommandRepository.UpdateUserPassword(ctx, request)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_UPDATE_PASSWORD", span, &status)
	}

	userResponse := s.mapper.ToUserResponse(res)

	keyId := fmt.Sprintf("user:id:%d", request.UserID)
	keyEmail := fmt.Sprintf("user:email:%s", res.Email)

	s.cacheCommand.DeleteCachedUser(ctx, keyId)
	s.cacheCommand.DeleteCachedUser(ctx, keyEmail)
	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)

	logSuccess("Successfully updated user password", zap.Int("user_id", res.ID))

	return userResponse, nil
}

func (s *userCommandService) TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", user_id))
	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.TrashedUser(ctx, user_id)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleTrashedUserError(err, method, "FAILED_TRASH_USER", span, &status)
	}

	userResponse := s.mapper.ToUserResponseDeleteAt(res)

	keyId := fmt.Sprintf("user:id:%d", user_id)
	keyEmail := fmt.Sprintf("user:email:%s", res.Email)

	s.cacheCommand.DeleteCachedUser(ctx, keyId)
	s.cacheCommand.DeleteCachedUser(ctx, keyEmail)
	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)
	s.cacheCommand.InvalidateTrashedUsers(ctx)

	logSuccess("Successfully trashed user", zap.Int("user_id", res.ID))

	return userResponse, nil
}

func (s *userCommandService) RestoreUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", user_id))
	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.RestoreUser(ctx, user_id)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleRestoreUserError(err, method, "FAILED_RESTORE_USER", span, &status)
	}

	userResponse := s.mapper.ToUserResponseDeleteAt(res)

	key := fmt.Sprintf("user:id:%d", user_id)

	keyEmail := fmt.Sprintf("user:email:%s", res.Email)

	s.cacheCommand.DeleteCachedUser(ctx, key)
	s.cacheCommand.DeleteCachedUser(ctx, keyEmail)
	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)
	s.cacheCommand.InvalidateTrashedUsers(ctx)

	logSuccess("Successfully restored user", zap.Int("user_id", res.ID))

	return userResponse, nil
}

func (s *userCommandService) DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", user_id))
	defer func() {
		end(status)
	}()

	_, err := s.userCommandRepository.DeleteUserPermanent(ctx, user_id)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleDeleteUserPermanentError(err, method, "FAILED_DELETE_PERMANENT", span, &status)
	}

	key := fmt.Sprintf("user:id:%d", user_id)
	keyEmail := "user:email:*"

	s.cacheCommand.DeleteCachedUser(ctx, key)
	s.cacheCommand.DeleteCachedUser(ctx, keyEmail)

	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)
	s.cacheCommand.InvalidateTrashedUsers(ctx)

	logSuccess("Successfully permanently deleted user", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userCommandService) RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.userCommandRepository.RestoreAllUser(ctx)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleRestoreAllUserError(err, method, "FAILED_RESTORE_ALL", span, &status)
	}

	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)
	s.cacheCommand.InvalidateTrashedUsers(ctx)

	logSuccess("Successfully restored all trashed users")

	return true, nil
}

func (s *userCommandService) DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.userCommandRepository.DeleteAllUserPermanent(ctx)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleDeleteAllUserPermanentError(err, method, "FAILED_DELETE_ALL_PERMANENT", span, &status)
	}

	s.cacheCommand.InvalidateAllUsers(ctx)
	s.cacheCommand.InvalidateActiveUsers(ctx)
	s.cacheCommand.InvalidateTrashedUsers(ctx)

	logSuccess("Successfully permanently deleted all trashed users")

	return true, nil
}
