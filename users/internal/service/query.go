package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	responsemapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"go.uber.org/zap"
)

type userQueryDeps struct {
	repository repository.UserQueryRepository
	logger     logger.LoggerInterface
	mapper     responsemapper.UserQueryResponseMapper
}

type userQueryService struct {
	userQueryRepository repository.UserQueryRepository
	logger              logger.LoggerInterface
	mapper              responsemapper.UserQueryResponseMapper
}

func NewUserQueryService(
	params *userQueryDeps,
) UserQueryService {
	return &userQueryService{
		userQueryRepository: params.repository,
		logger:              params.logger,
		mapper:              params.mapper,
	}
}

func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	if page != req.Page || pageSize != req.PageSize {
		s.logger.Warn("Pagination parameters normalized",
			zap.Int("original_page", req.Page),
			zap.Int("original_page_size", req.PageSize),
			zap.Int("normalized_page", page),
			zap.Int("normalized_page_size", pageSize),
		)
	}

	req.Page = page
	req.PageSize = pageSize

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve all users from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve users", 500)
	}

	usersResponse := s.mapper.ToUsersResponse(users)

	s.logger.Info("Successfully retrieved users",
		zap.Int("users_returned", len(usersResponse)),
		zap.Int("total_records", *totalRecords),
	)

	return usersResponse, totalRecords, nil
}

func (s *userQueryService) FindByID(ctx context.Context, id int) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("FindByID service called", zap.Int("user_id", id))

	user, err := s.userQueryRepository.FindById(ctx, id)

	if err != nil || user == nil {
		s.logger.Warn("User not found by ID", zap.Int("user_id", id), zap.Error(err))
		return nil, response.NewErrorResponse("user not found", 404)
	}

	userRes := s.mapper.ToUserResponse(user)

	s.logger.Info("Successfully found user by ID", zap.Int("user_id", id))
	return userRes, nil
}

func (s *userQueryService) FindByEmail(ctx context.Context, email string) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("FindByEmail service called", zap.String("email", email))

	user, err := s.userQueryRepository.FindByEmail(ctx, email)

	if err != nil || user == nil {
		s.logger.Warn("User not found by Email", zap.String("email", email), zap.Error(err))
		return nil, response.NewErrorResponse("user with the specified email not found", 404)
	}

	userRes := s.mapper.ToUserResponse(user)

	s.logger.Info("Successfully found user by Email", zap.String("email", email))
	return userRes, nil
}

func (s *userQueryService) FindByEmailAndVerify(ctx context.Context, email string) (*response.UserResponseWithPassword, *response.ErrorResponse) {
	s.logger.Info("FindByEmailAndVerify service called", zap.String("email", email))

	user, err := s.userQueryRepository.FindByEmailAndVerify(ctx, email)

	if err != nil || user == nil {
		s.logger.Warn("Authentication failed: invalid credentials for email", zap.String("email", email))
		return nil, response.NewErrorResponse("invalid credentials", 401)
	}

	userRes := s.mapper.ToUserWithPasswordResponse(user)

	s.logger.Info("Successfully found and verified user by Email", zap.String("email", email))
	return userRes, nil
}

func (s *userQueryService) FindByVerificationCode(ctx context.Context, code string) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("FindByVerificationCode service called", zap.String("code", code))

	user, err := s.userQueryRepository.FindByVerificationCode(ctx, code)

	if err != nil || user == nil {
		s.logger.Warn("Verification failed: invalid or expired code", zap.String("code", code))
		return nil, response.NewErrorResponse("invalid or expired verification code", 400)
	}

	userRes := s.mapper.ToUserResponse(user)

	s.logger.Info("Successfully found user by verification code", zap.String("code", code))
	return userRes, nil
}

func (s *userQueryService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	if page != req.Page || pageSize != req.PageSize {
		s.logger.Warn("Pagination parameters normalized for FindByActive",
			zap.Int("original_page", req.Page),
			zap.Int("normalized_page", page),
		)
	}

	req.Page = page
	req.PageSize = pageSize

	users, totalRecords, err := s.userQueryRepository.FindByActive(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve active users from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve active users", 500)
	}

	usersResponse := s.mapper.ToUsersResponseDeleteAt(users)

	s.logger.Info("Successfully retrieved active users",
		zap.Int("users_returned", len(usersResponse)),
		zap.Int("total_records", *totalRecords),
	)

	return usersResponse, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	if page != req.Page || pageSize != req.PageSize {
		s.logger.Warn("Pagination parameters normalized for FindByTrashed",
			zap.Int("original_page", req.Page),
			zap.Int("normalized_page", page),
		)
	}

	req.Page = page
	req.PageSize = pageSize

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed users from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve trashed users", 500)
	}

	usersResponse := s.mapper.ToUsersResponseDeleteAt(users)

	s.logger.Info("Successfully retrieved trashed users",
		zap.Int("users_returned", len(usersResponse)),
		zap.Int("total_records", *totalRecords),
	)

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
