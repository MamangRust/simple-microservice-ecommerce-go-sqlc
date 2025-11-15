package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
	userrepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors/user_errors/repository"
)

type userQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.UserQueryRecordMapper
}

func NewUserQueryRepository(db *db.Queries, mapper recordmapper.UserQueryRecordMapper) UserQueryRepository {
	return &userQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsers(ctx, reqDb)

	if err != nil {
		return nil, nil, userrepositoryerrors.ErrFindAllUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToUsersRecordPagination(res)

	return so, &totalCount, nil
}

func (r *userQueryRepository) FindById(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userrepositoryerrors.ErrUserNotFound
		}

		return nil, userrepositoryerrors.ErrUserNotFound
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, userrepositoryerrors.ErrUserNotFound
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

func (r *userQueryRepository) FindByEmailAndVerify(ctx context.Context, email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmailAndVerified(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userrepositoryerrors.ErrUserNotFound
		}

		return nil, userrepositoryerrors.ErrUserNotFound
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByVerificationCode(ctx, verification_code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userrepositoryerrors.ErrUserNotFound
		}

		return nil, userrepositoryerrors.ErrUserNotFound

	}

	return r.mapper.ToUserRecord(res), nil
}

func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersActive(ctx, reqDb)

	if err != nil {
		return nil, nil, userrepositoryerrors.ErrFindActiveUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToUsersRecordActivePagination(res)

	return so, &totalCount, nil
}

func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUserTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUserTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, userrepositoryerrors.ErrFindTrashedUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToUsersRecordTrashedPagination(res)

	return so, &totalCount, nil
}
