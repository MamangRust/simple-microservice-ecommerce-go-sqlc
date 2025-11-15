package rolerepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	rolerecordmapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/record/role"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
	rolerepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/repository"
)

type roleCommandRepository struct {
	db     *db.Queries
	mapper rolerecordmapper.RoleCommandRecordMapper
}

func NewRoleCommandRepository(db *db.Queries, mapper rolerecordmapper.RoleCommandRecordMapper) RoleCommandRepository {
	return &roleCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.CreateRole(ctx, req.Name)

	if err != nil {
		return nil, rolerepositoryerrors.ErrCreateRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID: int32(*req.ID),
		Name:   req.Name,
	})

	if err != nil {
		return nil, rolerepositoryerrors.ErrUpdateRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.TrashRole(ctx, int32(id))

	if err != nil {
		return nil, rolerepositoryerrors.ErrTrashedRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))

	if err != nil {
		return nil, rolerepositoryerrors.ErrRestoreRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(role_id))

	if err != nil {
		return false, rolerepositoryerrors.ErrDeleteRolePermanent
	}

	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)

	if err != nil {
		return false, rolerepositoryerrors.ErrRestoreAllRoles
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)

	if err != nil {
		return false, rolerepositoryerrors.ErrDeleteAllRoles
	}

	return true, nil
}
