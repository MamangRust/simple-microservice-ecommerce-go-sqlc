package rolehandler

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	roleprotomapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/proto/role"
	roleservice "github.com/MamangRust/simple_microservice_ecommerce/role/internal/service/role"
	rolegrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleCommandHandleGrpc struct {
	pbrole.UnimplementedRoleCommandServiceServer
	roleCommandService roleservice.RoleCommandService
	logger             logger.LoggerInterface
	mapper             roleprotomapper.RoleCommandProtoMapper
}

func NewRoleCommandHandleGrpc(command roleservice.RoleCommandService, logger logger.LoggerInterface) *roleCommandHandleGrpc {
	return &roleCommandHandleGrpc{
		roleCommandService: command,
		logger:             logger,
		mapper:             roleprotomapper.NewRoleCommandProtoMapper(),
	}
}

func (s *roleCommandHandleGrpc) CreateRole(ctx context.Context, request *pbrole.CreateRoleRequest) (*pbrole.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{
		Name: request.GetName(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("Invalid request data for CreateRole", zap.Error(err), zap.String("role_name", request.GetName()))
		return nil, rolegrpcerrors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommandService.CreateRole(ctx, req)

	if err != nil {
		s.logger.Error("Failed to create role in service",
			zap.String("error_message", err.Message),
			zap.String("role_name", request.GetName()),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Role successfully created",
		zap.Int("role_id", role.ID),
		zap.String("role_name", role.Name),
	)

	so := s.mapper.ToProtoResponseRole("success", "Successfully created role", role)
	return so, nil
}

func (s *roleCommandHandleGrpc) UpdateRole(ctx context.Context, request *pbrole.UpdateRoleRequest) (*pbrole.ApiResponseRole, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid RoleID received on update request", zap.Int("role_id", id))
		return nil, rolegrpcerrors.ErrGrpcRoleInvalidId
	}

	req := &requests.UpdateRoleRequest{
		ID:   &id,
		Name: request.GetName(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("Invalid request data for UpdateRole", zap.Error(err), zap.Int("role_id", id))
		return nil, rolegrpcerrors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommandService.UpdateRole(ctx, req)

	if err != nil {
		s.logger.Error("Failed to update role in service",
			zap.String("error_message", err.Message),
			zap.Int("role_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Role successfully updated", zap.Int("role_id", id))

	so := s.mapper.ToProtoResponseRole("success", "Successfully updated role", role)
	return so, nil
}

func (s *roleCommandHandleGrpc) TrashedRole(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRoleDeleteAt, error) {
	id := int(request.GetRoleId())

	if id == 0 {
		s.logger.Warn("Invalid RoleID received on trash request", zap.Int("role_id", id))
		return nil, rolegrpcerrors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommandService.TrashedRole(ctx, id)

	if err != nil {
		s.logger.Error("Failed to trash role in service",
			zap.String("error_message", err.Message),
			zap.Int("role_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Role successfully trashed", zap.Int("role_id", id))

	so := s.mapper.ToProtoResponseRoleDeleteAt("success", "Successfully trashed role", role)
	return so, nil
}

func (s *roleCommandHandleGrpc) RestoreRole(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRoleDeleteAt, error) {
	id := int(request.GetRoleId())

	if id == 0 {
		s.logger.Warn("Invalid RoleID received on restore request", zap.Int("role_id", id))
		return nil, rolegrpcerrors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommandService.RestoreRole(ctx, id)

	if err != nil {
		s.logger.Error("Failed to restore role in service",
			zap.String("error_message", err.Message),
			zap.Int("role_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Role successfully restored", zap.Int("role_id", id))

	so := s.mapper.ToProtoResponseRoleDeleteAt("success", "Successfully restored role", role)
	return so, nil
}

func (s *roleCommandHandleGrpc) DeleteRolePermanent(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRoleDelete, error) {
	id := int(request.GetRoleId())

	if id == 0 {
		s.logger.Warn("Invalid RoleID received on permanent delete request", zap.Int("role_id", id))
		return nil, rolegrpcerrors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommandService.DeleteRolePermanent(ctx, id)

	if err != nil {
		s.logger.Error("Failed to permanently delete role in service",
			zap.String("error_message", err.Message),
			zap.Int("role_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Role successfully permanently deleted", zap.Int("role_id", id))

	so := s.mapper.ToProtoResponseRoleDelete("success", "Successfully deleted role permanently")
	return so, nil
}

func (s *roleCommandHandleGrpc) RestoreAllRole(ctx context.Context, _ *emptypb.Empty) (*pbrole.ApiResponseRoleAll, error) {
	_, err := s.roleCommandService.RestoreAllRole(ctx)

	if err != nil {
		s.logger.Error("Failed to restore all roles in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	// Log: Informasi keberhasilan operasi massal
	s.logger.Info("All trashed roles successfully restored")

	so := s.mapper.ToProtoResponseRoleAll("success", "Successfully restored all roles")
	return so, nil
}

func (s *roleCommandHandleGrpc) DeleteAllRolePermanent(ctx context.Context, _ *emptypb.Empty) (*pbrole.ApiResponseRoleAll, error) {
	_, err := s.roleCommandService.DeleteAllRolePermanent(ctx)

	if err != nil {
		s.logger.Error("Failed to permanently delete all roles in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	// Log: Informasi keberhasilan operasi massal
	s.logger.Info("All trashed roles successfully permanently deleted")

	so := s.mapper.ToProtoResponseRoleAll("success", "Successfully deleted all roles permanently")
	return so, nil
}
