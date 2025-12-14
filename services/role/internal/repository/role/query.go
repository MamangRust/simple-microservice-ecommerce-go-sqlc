package rolerepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	rolerecordmapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/record/role"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
	rolerepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/repository"
)

type roleQueryRepository struct {
	db     *db.Queries
	mapper rolerecordmapper.RoleQueryRecordMapper
}

func NewRoleQueryRepository(db *db.Queries, mapper rolerecordmapper.RoleQueryRecordMapper) RoleQueryRepository {
	return &roleQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *roleQueryRepository) FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(ctx, reqDb)

	if err != nil {
		return nil, nil, rolerepositoryerrors.ErrFindAllRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToRolesRecordAll(res), &totalCount, nil
}

func (r *roleQueryRepository) FindById(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rolerepositoryerrors.ErrRoleNotFound
		}
		return nil, rolerepositoryerrors.ErrRoleNotFound
	}

	return r.mapper.ToRoleRecord(res), nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rolerepositoryerrors.ErrRoleNotFound
		}

		return nil, rolerepositoryerrors.ErrRoleNotFound
	}
	return r.mapper.ToRoleRecord(res), nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, user_id int) ([]*record.RoleRecord, error) {
	res, err := r.db.GetUserRoles(ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rolerepositoryerrors.ErrRoleNotFound
		}

		return nil, rolerepositoryerrors.ErrRoleNotFound
	}
	return r.mapper.ToRolesRecord(res), nil
}

func (r *roleQueryRepository) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(ctx, reqDb)

	if err != nil {
		return nil, nil, rolerepositoryerrors.ErrFindActiveRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToRolesRecordActive(res), &totalCount, nil
}

func (r *roleQueryRepository) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(ctx, reqDb)

	if err != nil {
		return nil, nil, rolerepositoryerrors.ErrFindTrashedRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToRolesRecordTrashed(res), &totalCount, nil
}
