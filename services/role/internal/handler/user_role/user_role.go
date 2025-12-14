package userrolehandler

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	userroleprotomapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/proto/user_role"
	userroleservice "github.com/MamangRust/simple_microservice_ecommerce/role/internal/service/user_role"
	userrolegrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/user_role_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
)

type userRoleHandleGrpc struct {
	pbuserrole.UnimplementedUserRoleServiceServer
	userRole userroleservice.UserRoleService
	logger   logger.LoggerInterface
	mapper   userroleprotomapper.UserRoleProtoMappper
}

func NewUserRoleHandleGrpc(userRole userroleservice.UserRoleService, logger logger.LoggerInterface) UserRoleHandleGrpc {
	return &userRoleHandleGrpc{
		userRole: userRole,
		logger:   logger,
		mapper:   userroleprotomapper.NewUserRoleProtoMapper(),
	}
}

func (s *userRoleHandleGrpc) AssignRole(ctx context.Context, request *pbuserrole.CreateUserRoleRequest) (*pbuserrole.ApiResponseUserRole, error) {
	userId := int(request.GetUserid())
	roleId := int(request.GetRoleid())

	if userId == 0 || roleId == 0 {
		return nil, userrolegrpcerrors.ErrGrpcRoleInvalidId
	}

	reqService := &requests.CreateUserRoleRequest{
		UserId: userId,
		RoleId: roleId,
	}

	userRole, err := s.userRole.AssignRoleToUser(ctx, reqService)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseUserRole("success", "Successfully assigned role to user", userRole)
	return so, nil
}

func (s *userRoleHandleGrpc) UpdateRole(ctx context.Context, request *pbuserrole.CreateUserRoleRequest) (*pbuserrole.ApiResponseUserRole, error) {
	userId := int(request.GetUserid())
	roleId := int(request.GetRoleid())

	if userId == 0 || roleId == 0 {
		return nil, userrolegrpcerrors.ErrGrpcRoleInvalidId
	}

	reqService := &requests.CreateUserRoleRequest{
		UserId: userId,
		RoleId: roleId,
	}

	userRole, err := s.userRole.UpdateRoleToUser(ctx, reqService)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseUserRole("success", "Successfully updated user role", userRole)
	return so, nil
}
