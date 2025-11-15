package roleservice

import (
	"context"
	"errors"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	roleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/role"
	repository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/role"
	rolerepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"go.uber.org/zap"
)

type roleCommandDeps struct {
	repository repository.RoleRepository
	logger     logger.LoggerInterface
	mapper     roleresponsemapper.RoleCommandResponseMapper
}

type roleCommandService struct {
	repository repository.RoleRepository
	logger     logger.LoggerInterface
	mapper     roleresponsemapper.RoleCommandResponseMapper
}

func NewRoleCommandService(deps *roleCommandDeps) RoleCommandService {
	return &roleCommandService{
		repository: deps.repository,
		logger:     deps.logger,
		mapper:     deps.mapper,
	}
}

func (s *roleCommandService) CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Info("Attempting to create role", zap.String("role_name", request.Name))

	existingRole, err := s.repository.FindByName(ctx, request.Name)
	if err != nil {
		s.logger.Error("Failed to check for existing role during creation", zap.Error(err), zap.String("role_name", request.Name))
		return nil, response.NewErrorResponse("failed to validate new role", 500)
	}
	if existingRole != nil {
		s.logger.Warn("Role creation failed: name already exists", zap.String("role_name", request.Name))
		return nil, response.NewErrorResponse("role with this name already exists", 409)
	}

	newRole, err := s.repository.CreateRole(ctx, request)
	if err != nil {
		s.logger.Error("Failed to create role in repository", zap.Error(err), zap.String("role_name", request.Name))
		return nil, response.NewErrorResponse("failed to create role", 500)
	}

	s.logger.Info("Role created successfully", zap.Int("role_id", newRole.ID), zap.String("role_name", newRole.Name))
	so := s.mapper.ToRoleResponse(newRole)
	return so, nil
}

func (s *roleCommandService) UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Info("Attempting to update role", zap.Int("role_id", *request.ID))

	existingRole, err := s.repository.FindById(ctx, *request.ID)
	if err != nil {
		s.logger.Error("Failed to find role for update", zap.Error(err), zap.Int("role_id", *request.ID))
		return nil, response.NewErrorResponse("role not found", 404)
	}

	if request.Name != "" && request.Name != existingRole.Name {
		s.logger.Info("Checking for duplicate role name during update", zap.String("new_name", request.Name))
		duplicateRole, err := s.repository.FindByName(ctx, request.Name)
		if err != nil && !errors.Is(err, rolerepositoryerrors.ErrRoleNotFound) {
			s.logger.Error("Failed to check for duplicate role name", zap.Error(err), zap.String("new_name", request.Name))
			return nil, response.NewErrorResponse("failed to check for duplicate role name", 500)
		}
		if duplicateRole != nil {
			s.logger.Warn("Role update failed: new name already exists", zap.String("new_name", request.Name))
			return nil, response.NewErrorResponse("role with this name already exists", 409)
		}
	}

	updatedRole, err := s.repository.UpdateRole(ctx, request)
	if err != nil {
		s.logger.Error("Failed to update role in repository", zap.Error(err), zap.Int("role_id", *request.ID))
		return nil, response.NewErrorResponse("failed to update role", 500)
	}

	s.logger.Info("Role updated successfully", zap.Int("role_id", updatedRole.ID), zap.String("role_name", updatedRole.Name))
	so := s.mapper.ToRoleResponse(updatedRole)
	return so, nil
}

func (s *roleCommandService) TrashedRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Info("Attempting to trash role", zap.Int("role_id", role_id))

	res, err := s.repository.TrashedRole(ctx, role_id)
	if err != nil {
		s.logger.Error("Failed to trash role in repository", zap.Error(err), zap.Int("role_id", role_id))
		return nil, response.NewErrorResponse("failed to trash role", 500)
	}

	s.logger.Info("Role trashed successfully", zap.Int("role_id", res.ID))
	so := s.mapper.ToRoleResponseDeleteAt(res)
	return so, nil
}

func (s *roleCommandService) RestoreRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Info("Attempting to restore role", zap.Int("role_id", role_id))

	res, err := s.repository.RestoreRole(ctx, role_id)
	if err != nil {
		s.logger.Error("Failed to restore role in repository", zap.Error(err), zap.Int("role_id", role_id))
		return nil, response.NewErrorResponse("failed to restore role", 500)
	}

	s.logger.Info("Role restored successfully", zap.Int("role_id", res.ID))
	so := s.mapper.ToRoleResponseDeleteAt(res)
	return so, nil
}

func (s *roleCommandService) DeleteRolePermanent(ctx context.Context, role_id int) (bool, *response.ErrorResponse) {
	s.logger.Info("Attempting to permanently delete role", zap.Int("role_id", role_id))

	_, err := s.repository.DeleteRolePermanent(ctx, role_id)
	if err != nil {
		s.logger.Error("Failed to permanently delete role in repository", zap.Error(err), zap.Int("role_id", role_id))
		return false, response.NewErrorResponse("failed to delete role permanently", 500)
	}

	s.logger.Info("Role permanently deleted successfully", zap.Int("role_id", role_id))
	return true, nil
}

func (s *roleCommandService) RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse) {
	s.logger.Info("Attempting to restore all trashed roles")

	_, err := s.repository.RestoreAllRole(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all roles in repository", zap.Error(err))
		return false, response.NewErrorResponse("failed to restore all roles", 500)
	}

	s.logger.Info("All trashed roles restored successfully")
	return true, nil
}

func (s *roleCommandService) DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	s.logger.Info("Attempting to permanently delete all trashed roles")

	_, err := s.repository.DeleteAllRolePermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all roles in repository", zap.Error(err))
		return false, response.NewErrorResponse("failed to delete all roles permanently", 500)
	}

	s.logger.Info("All trashed roles permanently deleted successfully")
	return true, nil
}
