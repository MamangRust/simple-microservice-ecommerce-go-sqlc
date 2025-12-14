package userroleservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/errorhandler"
	userroleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/user_role"
	userrolerepository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/user_role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type UserRoleServiceDeps struct {
	Repository   userrolerepository.UserRoleRepository
	Errorhandler errorhandler.UserRoleErrorHandler
	Logger       logger.LoggerInterface
}

type userRoleService struct {
	repository    userrolerepository.UserRoleRepository
	logger        logger.LoggerInterface
	mapper        userroleresponsemapper.UserRoleResponseMappping
	observability observability.TraceLoggerObservability
	errorhandler  errorhandler.UserRoleErrorHandler
}

func NewUserRoleService(params *UserRoleServiceDeps) UserRoleService {
	observability, _ := observability.NewObservability("user-role-service", params.Logger)

	return &userRoleService{
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        userroleresponsemapper.NewUserRoleResponseMapper(),
		errorhandler:  params.Errorhandler,
		observability: observability,
	}
}

func (s *userRoleService) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.UserRoleResponse, *response.ErrorResponse) {
	const method = "AssignRoleToUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", req.UserId), attribute.Int("role_id", req.RoleId))

	defer func() {
		end(status)
	}()

	res, err := s.repository.AssignRoleToUser(ctx, req)

	if err != nil {
		status = "error"

		return s.errorhandler.HandleAssignRoleToUserError(err, method, "FAILED_ASSIGN_ROLE_TO_USER", span, &status, zap.Int("role_id", req.RoleId), zap.Int("user_id", req.UserId))
	}

	so := s.mapper.ToUserRoleResponse(res)

	logSuccess("Successfully assigned role to user",
		zap.Int("user_id", req.UserId),
		zap.Int("role_id", req.RoleId))

	return so, nil
}

func (s *userRoleService) UpdateRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.UserRoleResponse, *response.ErrorResponse) {
	const method = "UpdateRoleToUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", req.UserId), attribute.Int("role_id", req.RoleId))

	defer func() {
		end(status)
	}()

	res, err := s.repository.UpdateRoleToUser(ctx, req)

	if err != nil {
		status = "error"

		return s.errorhandler.HandleUpdateRoleToUserError(err, method, "FAILED_UPDATE_ROLE_TO_USER", span, &status, zap.Error(err),
			zap.Int("user_id", req.UserId),
			zap.Int("role_id", req.RoleId))
	}

	logSuccess("Successfully updated role for user",
		zap.Int("user_id", req.UserId),
		zap.Int("role_id", req.RoleId),
	)

	so := s.mapper.ToUserRoleResponse(res)
	return so, nil
}
