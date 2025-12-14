package handler

import (
	"context"
	"math"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	protomapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/proto"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/service"
	usergrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors/user_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type userQueryHandleGrpc struct {
	pbuser.UnimplementedUserQueryServiceServer
	userQueryService service.UserQueryService
	logger           logger.LoggerInterface
	mapper           protomapper.UserQueryProtoMapper
}

func NewUserQueryHandleGrpc(query service.UserQueryService, logger logger.LoggerInterface) pbuser.UserQueryServiceServer {
	return &userQueryHandleGrpc{
		userQueryService: query,
		logger:           logger,
		mapper:           protomapper.NewUserQueryProtoMapper(),
	}
}

func (s *userQueryHandleGrpc) FindAll(ctx context.Context, request *pbuser.FindAllUserRequest) (*pbuser.ApiResponsePaginationUser, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(ctx, reqService)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	Pagination := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationUser(Pagination, "success", "Successfully fetched users", users)
	return so, nil
}

func (s *userQueryHandleGrpc) FindByActive(ctx context.Context, request *pbuser.FindAllUserRequest) (*pbuser.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(ctx, reqService)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	Pagination := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationUserDeleteAt(Pagination, "success", "Successfully fetched active users", users)

	return so, nil
}

func (s *userQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pbuser.FindAllUserRequest) (*pbuser.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByTrashed(ctx, reqService)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	Pagination := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationUserDeleteAt(Pagination, "success", "Successfully fetched trashed users", users)

	return so, nil
}

func (s *userQueryHandleGrpc) FindById(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUser, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, usergrpcerrors.ErrGrpcUserInvalidId
	}

	user, err := s.userQueryService.FindByID(ctx, id)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully fetched user", user)

	return so, nil
}

func (s *userQueryHandleGrpc) FindVerificationCode(ctx context.Context, request *pbuser.VerifyCodeRequest) (*pbuser.ApiResponseUser, error) {
	user, err := s.userQueryService.FindByVerificationCode(ctx, request.Code)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully fetched user", user)

	return so, nil
}

func (s *userQueryHandleGrpc) FindByEmail(ctx context.Context, request *pbuser.FindByEmailUserRequest) (*pbuser.ApiResponseUser, error) {
	user, err := s.userQueryService.FindByEmail(ctx, request.Email)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully fetched user", user)

	return so, nil
}

func (s *userQueryHandleGrpc) FindByEmailAndVerify(ctx context.Context, request *pbuser.FindByEmailUserRequest) (*pbuser.ApiResponseUserWithPassword, error) {
	user, err := s.userQueryService.FindByEmailAndVerify(ctx, request.Email)

	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseUserWithPassword("success", "Successfully fetched user", user)

	return so, nil
}
