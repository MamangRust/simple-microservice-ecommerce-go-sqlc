package userrolerepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	userrolerecordmapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/record/user_role"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
	userrolerepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/user_role_errors/repository"
)

type userRoleRepository struct {
	db     *db.Queries
	mapper userrolerecordmapper.UserRoleRecordMapping
}

func NewUserRoleRepository(db *db.Queries, mapper userrolerecordmapper.UserRoleRecordMapping) UserRoleRepository {
	return &userRoleRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error) {
	res, err := r.db.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return nil, userrolerepositoryerrors.ErrAssignRoleToUser
	}

	return r.mapper.ToUserRoleRecord(res), nil
}

func (r *userRoleRepository) UpdateRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error) {
	res, err := r.db.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return nil, userrolerepositoryerrors.ErrUpdateRoleToUser
	}

	return r.mapper.ToUserRoleRecord(res), nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.db.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return userrolerepositoryerrors.ErrRemoveRole
	}

	return nil
}
