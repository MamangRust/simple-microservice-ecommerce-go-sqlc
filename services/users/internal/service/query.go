package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/errorhandler"
	responsemapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/service"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/user/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userQueryDeps struct {
	repository   repository.UserQueryRepository
	logger       logger.LoggerInterface
	mapper       responsemapper.UserQueryResponseMapper
	erorrhandler errorhandler.UserQueryErrorHandler
	mencache     mencache.UserQueryCache
}

type userQueryService struct {
	userQueryRepository repository.UserQueryRepository
	logger              logger.LoggerInterface
	mapper              responsemapper.UserQueryResponseMapper
	observability       observability.TraceLoggerObservability
	errorhandler        errorhandler.UserQueryErrorHandler
	mencache            mencache.UserQueryCache
}

func NewUserQueryService(
	params *userQueryDeps,
) UserQueryService {
	observability, _ := observability.NewObservability("user-query-service", params.logger)

	return &userQueryService{
		userQueryRepository: params.repository,
		logger:              params.logger,
		mapper:              params.mapper,
		observability:       observability,
		errorhandler:        params.erorrhandler,
		mencache:            params.mencache,
	}
}

func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUsers(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(ctx, req)

	if err != nil {
		status = "error"

		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_USERS", span, &status, zap.String("error", err.Error()))
	}

	usersResponse := s.mapper.ToUsersResponse(users)

	s.mencache.SetCachedUsers(ctx, req, usersResponse, totalRecords)

	logSuccess("Successfully retrieved users", zap.Int("users_returned", len(usersResponse)))

	return usersResponse, totalRecords, nil
}

func (s *userQueryService) FindByID(ctx context.Context, id int) (*response.UserResponse, *response.ErrorResponse) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", id))

	defer func() {
		end(status)
	}()

	user, err := s.userQueryRepository.FindById(ctx, id)

	defaultErr := response.NewErrorResponse("user not found", 404)
	if err != nil || user == nil {
		status = "error"

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_ID", span, &status, defaultErr, zap.Int("user_id", id))
	}

	userRes := s.mapper.ToUserResponse(user)

	s.mencache.SetCachedUserById(ctx, userRes)

	logSuccess("Successfully found user by ID", zap.Int("user_id", id))

	return userRes, nil
}

func (s *userQueryService) FindByEmail(ctx context.Context, email string) (*response.UserResponse, *response.ErrorResponse) {
	const method = "FindByEmail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", email))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedUserByEmail(ctx, email); found {
		logSuccess("Data found in cache", zap.String("user.email", email))

		return data, nil
	}

	user, err := s.userQueryRepository.FindByEmail(ctx, email)

	defaultErr := response.NewErrorResponse("user with the specified email not found", 404)
	if err != nil || user == nil {
		status = "error"

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_EMAIL", span, &status, defaultErr, zap.String("email", email))
	}

	userRes := s.mapper.ToUserResponse(user)

	s.mencache.SetCachedUserByEmail(ctx, userRes)

	logSuccess("Successfully found user by Email", zap.String("email", email))

	return userRes, nil
}

func (s *userQueryService) FindByEmailAndVerify(ctx context.Context, email string) (*response.UserResponseWithPassword, *response.ErrorResponse) {
	const method = "FindByEmailAndVerify"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", email))

	defer func() {
		end(status)
	}()

	user, err := s.userQueryRepository.FindByEmailAndVerify(ctx, email)

	defaultErr := response.NewErrorResponse("invalid credentials", 401)
	if err != nil || user == nil {
		status = "error"

		return s.errorhandler.HandleRepositorySingleWithPasswordError(err, method, "FAILED_VERIFY_BY_EMAIL", span, &status, defaultErr, zap.String("email", email))
	}

	userRes := s.mapper.ToUserWithPasswordResponse(user)

	logSuccess("Successfully found and verified user by Email", zap.String("email", email))

	return userRes, nil
}

func (s *userQueryService) FindByVerificationCode(ctx context.Context, code string) (*response.UserResponse, *response.ErrorResponse) {
	const method = "FindByVerificationCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("code", code))

	defer func() {
		end(status)
	}()

	user, err := s.userQueryRepository.FindByVerificationCode(ctx, code)

	defaultErr := response.NewErrorResponse("invalid or expired verification code", 400)
	if err != nil || user == nil {
		status = "error"

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_CODE", span, &status, defaultErr, zap.String("code", code))
	}

	userRes := s.mapper.ToUserResponse(user)

	logSuccess("Successfully found user by verification code", zap.String("code", code))

	return userRes, nil
}

func (s *userQueryService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"
	search := req.Search
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByActive(ctx, req)

	defaultErr := response.NewErrorResponse("failed to retrieve active users", 500)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_ACTIVE", span, &status, defaultErr)
	}

	usersResponse := s.mapper.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserActive(ctx, req, usersResponse, totalRecords)

	logSuccess("Successfully retrieved active users", zap.Int("users_returned", len(usersResponse)), zap.Int("total_records", *totalRecords))

	return usersResponse, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"
	search := req.Search
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(ctx, req)

	defaultErr := response.NewErrorResponse("failed to retrieve trashed users", 500)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_TRASHED", span, &status, defaultErr)
	}

	usersResponse := s.mapper.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserTrashed(ctx, req, usersResponse, totalRecords)

	logSuccess("Successfully retrieved trashed users", zap.Int("users_returned", len(usersResponse)), zap.Int("total_records", *totalRecords))

	return usersResponse, totalRecords, nil
}

func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
