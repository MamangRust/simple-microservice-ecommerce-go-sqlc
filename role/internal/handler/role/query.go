package rolehandler

import (
	"context"
	"math"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	roleprotomapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/proto/role"
	roleservice "github.com/MamangRust/simple_microservice_ecommerce/role/internal/service/role"
	rolegrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type roleQueryHandleGrpc struct {
	pbrole.UnimplementedRoleQueryServiceServer
	roleQueryService roleservice.RoleQueryService
	logger           logger.LoggerInterface
	mapper           roleprotomapper.RoleQueryProtoMapper
}

func NewRoleQueryHandleGrpc(query roleservice.RoleQueryService, logger logger.LoggerInterface) RoleQueryHandleGrpc {
	return &roleQueryHandleGrpc{
		roleQueryService: query,
		logger:           logger,
		mapper:           roleprotomapper.NewRoleQueryProtoMapper(),
	}
}

func (s *roleQueryHandleGrpc) FindAllRole(ctx context.Context, request *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRole, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		s.logger.Warn("Invalid page number received, defaulting to 1", zap.Int("page", page))
		page = 1
	}
	if pageSize <= 0 {
		s.logger.Warn("Invalid pageSize received, defaulting to 10", zap.Int("pageSize", pageSize))
		pageSize = 10
	}

	reqService := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQueryService.FindAll(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to find all roles", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched roles",
		zap.Int("count", len(roles)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationRole(paginationMeta, "success", "Successfully fetched roles", roles)
	return so, nil
}

func (s *roleQueryHandleGrpc) FindByIdRole(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	id := int(request.GetRoleId())

	if id == 0 {
		s.logger.Warn("Invalid RoleID received", zap.Int("roleId", id))
		return nil, rolegrpcerrors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQueryService.FindById(ctx, id)

	if err != nil {
		s.logger.Error("Failed to find role by ID", zap.String("error_message", err.Message), zap.Int("roleId", id))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Successfully fetched role", zap.Int("roleId", id))

	so := s.mapper.ToProtoResponseRole("success", "Successfully fetched role", role)
	return so, nil
}

func (s *roleQueryHandleGrpc) FindByName(ctx context.Context, request *pbrole.FindByNameRequest) (*pbrole.ApiResponseRole, error) {
	name := request.GetName()

	if name == "" {
		s.logger.Warn("Invalid RoleName received (cannot be empty)", zap.String("roleName", name))
		return nil, status.Error(codes.InvalidArgument, "role name cannot be empty")
	}

	role, err := s.roleQueryService.FindByName(ctx, name)

	if err != nil {
		s.logger.Error("Failed to find role by name", zap.String("error_message", err.Message), zap.String("roleName", name))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Successfully fetched role by name", zap.String("roleName", name))

	so := s.mapper.ToProtoResponseRole("success", "Successfully fetched role by name", role)
	return so, nil
}

func (s *roleQueryHandleGrpc) FindByActive(ctx context.Context, request *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		s.logger.Warn("Invalid page number received, defaulting to 1", zap.Int("page", page))
		page = 1
	}
	if pageSize <= 0 {
		s.logger.Warn("Invalid pageSize received, defaulting to 10", zap.Int("pageSize", pageSize))
		pageSize = 10
	}

	reqService := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQueryService.FindByActiveRole(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to find active roles", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active roles",
		zap.Int("count", len(roles)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched active roles", roles)
	return so, nil
}

func (s *roleQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		s.logger.Warn("Invalid page number received, defaulting to 1", zap.Int("page", page))
		page = 1
	}
	if pageSize <= 0 {
		s.logger.Warn("Invalid pageSize received, defaulting to 10", zap.Int("pageSize", pageSize))
		pageSize = 10
	}

	reqService := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQueryService.FindByTrashedRole(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to find trashed roles", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed roles",
		zap.Int("count", len(roles)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched trashed roles", roles)
	return so, nil
}

func (s *roleQueryHandleGrpc) FindByUserId(ctx context.Context, request *pbrole.FindByIdUserRoleRequest) (*pbrole.ApiResponsesRole, error) {
	userId := int(request.GetUserId())

	if userId == 0 {
		s.logger.Warn("Invalid UserID received", zap.Int("userId", userId))
		return nil, rolegrpcerrors.ErrGrpcRoleInvalidId
	}

	roles, err := s.roleQueryService.FindByUserId(ctx, userId)

	if err != nil {
		s.logger.Error("Failed to find roles by UserID", zap.String("error_message", err.Message), zap.Int("userId", userId))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Successfully fetched roles for user",
		zap.Int("userId", userId),
		zap.Int("count", len(roles)),
	)

	so := s.mapper.ToProtoResponsesRole("success", "Successfully fetched roles for user", roles)
	return so, nil
}
