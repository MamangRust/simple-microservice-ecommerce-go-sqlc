package userroleservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	userroleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/user_role"
	userrolerepository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/user_role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"go.uber.org/zap"
)

type UserRoleServiceDeps struct {
	Repository userrolerepository.UserRoleRepository
}

type userRoleService struct {
	repository userrolerepository.UserRoleRepository
	logger     logger.LoggerInterface
	mapper     userroleresponsemapper.UserRoleResponseMappping
}

func NewUserRoleService(repository userrolerepository.UserRoleRepository, logger logger.LoggerInterface) UserRoleService {
	return &userRoleService{
		repository: repository,
		logger:     logger,
		mapper:     userroleresponsemapper.NewUserRoleResponseMapper(),
	}
}

func (s *userRoleService) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.UserRoleResponse, *response.ErrorResponse) {
	res, err := s.repository.AssignRoleToUser(ctx, req)

	if err != nil {
		s.logger.Error("Failed to assign role to user in repository",
			zap.Error(err),
			zap.Int("user_id", req.UserId),
			zap.Int("role_id", req.RoleId),
		)
		return nil, response.NewErrorResponse("failed to assign role to user", 409)
	}

	s.logger.Info("Successfully assigned role to user",
		zap.Int("user_id", req.UserId),
		zap.Int("role_id", req.RoleId),
	)

	so := s.mapper.ToUserRoleResponse(res)
	return so, nil
}

func (s *userRoleService) UpdateRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.UserRoleResponse, *response.ErrorResponse) {
	res, err := s.repository.UpdateRoleToUser(ctx, req)

	if err != nil {
		s.logger.Error("Failed to update role for user in repository",
			zap.Error(err),
			zap.Int("user_id", req.UserId),
			zap.Int("role_id", req.RoleId),
		)
		return nil, response.NewErrorResponse("user role assignment not found", 409)
	}

	s.logger.Info("Successfully updated role for user",
		zap.Int("user_id", req.UserId),
		zap.Int("role_id", req.RoleId),
	)

	so := s.mapper.ToUserRoleResponse(res)
	return so, nil
}
