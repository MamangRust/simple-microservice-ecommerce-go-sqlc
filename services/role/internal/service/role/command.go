package roleservice

import (
	"context"
	"errors"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/errorhandler"
	roleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/role"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/role/internal/redis"
	repository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/role"
	rolerepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type roleCommandDeps struct {
	repository   repository.RoleRepository
	logger       logger.LoggerInterface
	mapper       roleresponsemapper.RoleCommandResponseMapper
	errorhandler errorhandler.RoleCommandErrorHandler
	cacheQuery   mencache.RoleQueryCache
	cacheCommand mencache.RoleCommandCache
}

type roleCommandService struct {
	repository    repository.RoleRepository
	logger        logger.LoggerInterface
	mapper        roleresponsemapper.RoleCommandResponseMapper
	errorhandler  errorhandler.RoleCommandErrorHandler
	cacheQuery    mencache.RoleQueryCache
	cacheCommand  mencache.RoleCommandCache
	observability observability.TraceLoggerObservability
}

func NewRoleCommandService(deps *roleCommandDeps) RoleCommandService {
	observability, _ := observability.NewObservability("role-command-service", deps.logger)

	return &roleCommandService{
		repository:    deps.repository,
		logger:        deps.logger,
		mapper:        deps.mapper,
		errorhandler:  deps.errorhandler,
		cacheQuery:    deps.cacheQuery,
		cacheCommand:  deps.cacheCommand,
		observability: observability,
	}
}

func (s *roleCommandService) CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "CreateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("role_name", request.Name))
	defer func() {
		end(status)
	}()

	existingRole, err := s.repository.FindByName(ctx, request.Name)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleCreateRoleError(err, method, "FAILED_CHECK_EXISTING_ROLE", span, &status, zap.String("role_name", request.Name))
	}
	if existingRole != nil {
		status = "error"
		customErr := errors.New("role with this name already exists")
		return s.errorhandler.HandleCreateRoleError(customErr, method, "ROLE_ALREADY_EXISTS", span, &status, zap.String("role_name", request.Name))
	}

	newRole, err := s.repository.CreateRole(ctx, request)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleCreateRoleError(err, method, "FAILED_CREATE_ROLE", span, &status, zap.String("role_name", request.Name))
	}

	so := s.mapper.ToRoleResponse(newRole)

	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)

	logSuccess("Successfully created role", zap.Int("role_id", newRole.ID), zap.String("role_name", newRole.Name))

	return so, nil
}

func (s *roleCommandService) UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "UpdateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("role_id", *request.ID))
	defer func() {
		end(status)
	}()

	existingRole, err := s.repository.FindById(ctx, *request.ID)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateRoleError(err, method, "FAILED_FIND_ROLE", span, &status, zap.Int("role_id", *request.ID))
	}

	if request.Name != "" && request.Name != existingRole.Name {
		duplicateRole, err := s.repository.FindByName(ctx, request.Name)
		if err != nil && !errors.Is(err, rolerepositoryerrors.ErrRoleNotFound) {
			status = "error"
			return s.errorhandler.HandleUpdateRoleError(err, method, "FAILED_CHECK_DUPLICATE_NAME", span, &status, zap.String("new_name", request.Name))
		}
		if duplicateRole != nil {
			status = "error"
			customErr := errors.New("role with this name already exists")
			return s.errorhandler.HandleUpdateRoleError(customErr, method, "ROLE_NAME_ALREADY_EXISTS", span, &status, zap.String("new_name", request.Name))
		}
	}

	updatedRole, err := s.repository.UpdateRole(ctx, request)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateRoleError(err, method, "FAILED_UPDATE_ROLE", span, &status, zap.Int("role_id", *request.ID))
	}

	so := s.mapper.ToRoleResponse(updatedRole)

	s.cacheCommand.DeleteCachedRole(ctx, *request.ID)
	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)

	logSuccess("Successfully updated role", zap.Int("role_id", updatedRole.ID), zap.String("role_name", updatedRole.Name))

	return so, nil
}

func (s *roleCommandService) TrashedRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("role_id", role_id))
	defer func() {
		end(status)
	}()

	res, err := s.repository.TrashedRole(ctx, role_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleTrashedRoleError(err, method, "FAILED_TRASH_ROLE", span, &status, zap.Int("role_id", role_id))
	}

	so := s.mapper.ToRoleResponseDeleteAt(res)

	s.cacheCommand.DeleteCachedRole(ctx, role_id)
	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)
	s.cacheCommand.InvalidateTrashedRoles(ctx)

	logSuccess("Successfully trashed role", zap.Int("role_id", res.ID))

	return so, nil
}

func (s *roleCommandService) RestoreRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("role_id", role_id))
	defer func() {
		end(status)
	}()

	res, err := s.repository.RestoreRole(ctx, role_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleRestoreRoleError(err, method, "FAILED_RESTORE_ROLE", span, &status, zap.Int("role_id", role_id))
	}

	so := s.mapper.ToRoleResponseDeleteAt(res)

	s.cacheCommand.DeleteCachedRole(ctx, role_id)
	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)
	s.cacheCommand.InvalidateTrashedRoles(ctx)

	logSuccess("Successfully restored role", zap.Int("role_id", res.ID))

	return so, nil
}

func (s *roleCommandService) DeleteRolePermanent(ctx context.Context, role_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("role_id", role_id))
	defer func() {
		end(status)
	}()

	_, err := s.repository.DeleteRolePermanent(ctx, role_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleDeleteRolePermanentError(err, method, "FAILED_DELETE_ROLE_PERMANENT", span, &status, zap.Int("role_id", role_id))
	}

	s.cacheCommand.DeleteCachedRole(ctx, role_id)
	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)
	s.cacheCommand.InvalidateTrashedRoles(ctx)

	logSuccess("Successfully permanently deleted role", zap.Int("role_id", role_id))

	return true, nil
}

func (s *roleCommandService) RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.repository.RestoreAllRole(ctx)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleRestoreAllRoleError(err, method, "FAILED_RESTORE_ALL_ROLES", span, &status)
	}

	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)
	s.cacheCommand.InvalidateTrashedRoles(ctx)

	logSuccess("Successfully restored all trashed roles")

	return true, nil
}

func (s *roleCommandService) DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.repository.DeleteAllRolePermanent(ctx)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleDeleteAllRolePermanentError(err, method, "FAILED_DELETE_ALL_ROLES_PERMANENT", span, &status)
	}

	s.cacheCommand.InvalidateAllRoles(ctx)
	s.cacheCommand.InvalidateActiveRoles(ctx)
	s.cacheCommand.InvalidateTrashedRoles(ctx)

	logSuccess("Successfully permanently deleted all trashed roles")

	return true, nil
}
