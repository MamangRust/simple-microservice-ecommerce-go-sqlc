package roleservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/errorhandler"
	roleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/role"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/role/internal/redis"
	repository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type roleQueryDeps struct {
	repository   repository.RoleQueryRepository
	logger       logger.LoggerInterface
	mapper       roleresponsemapper.RoleQueryResponseMapper
	mencache     mencache.RoleQueryCache
	errorhandler errorhandler.RoleQueryErrorHandler
}

type roleQueryService struct {
	repository    repository.RoleQueryRepository
	logger        logger.LoggerInterface
	mapper        roleresponsemapper.RoleQueryResponseMapper
	mencache      mencache.RoleQueryCache
	errorhandler  errorhandler.RoleQueryErrorHandler
	observability observability.TraceLoggerObservability
}

func NewRoleQueryService(deps *roleQueryDeps) RoleQueryService {
	observability, _ := observability.NewObservability("role-query-service", deps.logger)

	return &roleQueryService{
		repository:    deps.repository,
		logger:        deps.logger,
		mapper:        deps.mapper,
		mencache:      deps.mencache,
		errorhandler:  deps.errorhandler,
		observability: observability,
	}
}

func (s *roleQueryService) FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoles(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.repository.FindAllRoles(ctx, req)

	if err != nil {
		status = "error"
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_ROLES", span, &status)
	}

	so := s.mapper.ToRolesResponse(res)

	s.mencache.SetCachedRoles(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved roles",
		zap.Int("roles_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByName(ctx context.Context, name string) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindByName"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("name", name))

	defer func() {
		end(status)
	}()

	res, err := s.repository.FindByName(ctx, name)

	if err != nil || res == nil {
		status = "error"
		defaultErr := response.NewErrorResponse("role not found", 404)
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_NAME", span, &status, defaultErr,
			zap.String("role_name", name))
	}

	so := s.mapper.ToRoleResponse(res)

	s.mencache.SetCachedRoleByName(ctx, so)

	logSuccess("Successfully found role by name", zap.String("role_name", name))
	return so, nil
}

func (s *roleQueryService) FindById(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleById(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("id", id))
		return data, nil
	}

	res, err := s.repository.FindById(ctx, id)

	if err != nil || res == nil {
		status = "error"
		defaultErr := response.NewErrorResponse("role not found", 404)
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_ID", span, &status, defaultErr,
			zap.Int("role_id", id))
	}

	so := s.mapper.ToRoleResponse(res)

	s.mencache.SetCachedRoleById(ctx, so)

	logSuccess("Successfully found role by ID", zap.Int("role_id", id))
	return so, nil
}

func (s *roleQueryService) FindByUserId(ctx context.Context, userId int) ([]*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("userId", userId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleByUserId(ctx, userId); found {
		logSuccess("Data found in cache", zap.Int("userId", userId))
		return data, nil
	}

	res, err := s.repository.FindByUserId(ctx, userId)

	if err != nil || res == nil {
		status = "error"
		return s.errorhandler.HandleRepositoryListError(err, method, "FAILED_FIND_BY_USER_ID", span, &status,
			zap.Int("user_id", userId))
	}

	so := s.mapper.ToRolesResponse(res)

	s.mencache.SetCachedRoleByUserId(ctx, userId, so)

	logSuccess("Successfully found roles for user",
		zap.Int("user_id", userId),
		zap.Int("roles_count", len(so)),
	)
	return so, nil
}

func (s *roleQueryService) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindByActiveRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.repository.FindByActiveRole(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve active roles", 500)
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_ACTIVE_ROLES", span, &status, errResp)
	}

	so := s.mapper.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleActive(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved active roles",
		zap.Int("roles_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindByTrashedRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.repository.FindByTrashedRole(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve trashed roles", 500)
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_TRASHED_ROLES", span, &status, errResp)
	}

	so := s.mapper.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved trashed roles",
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
