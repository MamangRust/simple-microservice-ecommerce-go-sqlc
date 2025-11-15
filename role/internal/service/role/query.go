package roleservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	roleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/role"
	repository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"go.uber.org/zap"
)

type roleQueryDeps struct {
	repository repository.RoleQueryRepository
	logger     logger.LoggerInterface
	mapper     roleresponsemapper.RoleQueryResponseMapper
}

type roleQueryService struct {
	repository repository.RoleQueryRepository
	logger     logger.LoggerInterface
	mapper     roleresponsemapper.RoleQueryResponseMapper
}

func NewRoleQueryService(deps *roleQueryDeps) RoleQueryService {
	return &roleQueryService{
		repository: deps.repository,
		logger:     deps.logger,
		mapper:     deps.mapper,
	}
}

func (s *roleQueryService) FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	res, totalRecords, err := s.repository.FindAllRoles(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve all roles from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve roles", 500)
	}

	so := s.mapper.ToRolesResponse(res)

	s.logger.Info("Successfully retrieved roles",
		zap.Int("roles_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByName(ctx context.Context, name string) (*response.RoleResponse, *response.ErrorResponse) {
	res, err := s.repository.FindByName(ctx, name)

	if err != nil || res == nil {
		s.logger.Warn("Role not found by name", zap.String("role_name", name), zap.Error(err))
		return nil, response.NewErrorResponse("role not found", 404)
	}

	so := s.mapper.ToRoleResponse(res)

	s.logger.Info("Successfully found role by name", zap.String("role_name", name))
	return so, nil
}

func (s *roleQueryService) FindById(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	res, err := s.repository.FindById(ctx, id)

	if err != nil || res == nil {
		s.logger.Warn("Role not found by ID", zap.Int("role_id", id), zap.Error(err))
		return nil, response.NewErrorResponse("role not found", 404)
	}

	so := s.mapper.ToRoleResponse(res)

	s.logger.Info("Successfully found role by ID", zap.Int("role_id", id))
	return so, nil
}

func (s *roleQueryService) FindByUserId(ctx context.Context, userId int) ([]*response.RoleResponse, *response.ErrorResponse) {
	res, err := s.repository.FindByUserId(ctx, userId)

	if err != nil || res == nil {
		s.logger.Warn("Roles not found for user", zap.Int("user_id", userId), zap.Error(err))
		return nil, response.NewErrorResponse("roles not found for this user", 404)
	}

	so := s.mapper.ToRolesResponse(res)

	s.logger.Info("Successfully found roles for user",
		zap.Int("user_id", userId),
		zap.Int("roles_count", len(so)),
	)
	return so, nil
}

func (s *roleQueryService) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	res, totalRecords, err := s.repository.FindByActiveRole(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve active roles from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve active roles", 500)
	}

	so := s.mapper.ToRolesResponseDeleteAt(res)

	s.logger.Info("Successfully retrieved active roles",
		zap.Int("roles_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	res, totalRecords, err := s.repository.FindByTrashedRole(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed roles from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve trashed roles", 500)
	}

	so := s.mapper.ToRolesResponseDeleteAt(res)

	s.logger.Info("Successfully retrieved trashed roles",
		zap.Int("roles_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *roleQueryService) normalizePagination(page, pageSize int) (int, int) {
	originalPage, originalPageSize := page, pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if page != originalPage || pageSize != originalPageSize {
		s.logger.Warn("Pagination parameters normalized",
			zap.Int("original_page", originalPage),
			zap.Int("original_page_size", originalPageSize),
			zap.Int("new_page", page),
			zap.Int("new_page_size", pageSize),
		)
	}

	return page, pageSize
}
