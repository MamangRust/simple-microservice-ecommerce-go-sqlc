package repository

import (
	"context"
	"database/sql"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
	userrepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors/user_errors/repository"
)

type userCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.UserCommandRecordMapper
}

func NewUserCommandRepository(db *db.Queries, mapper recordmapper.UserCommandRecordMapper) UserCommandRepository {
	return &userCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       sql.NullBool{Bool: request.IsVerified, Valid: true},
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		return nil, userrepositoryerrors.ErrCreateUser
	}

	so := r.mapper.ToUserRecord(user)

	return so, nil
}

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*record.UserRecord, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(ctx, req)

	if err != nil {
		return nil, userrepositoryerrors.ErrUpdateUser
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

func (r *userCommandRepository) UpdateUserIsVerified(ctx context.Context, req *requests.UpdateUserVerifiedRequest) (*record.UserRecord, error) {
	res, err := r.db.UpdateUserIsVerified(ctx, db.UpdateUserIsVerifiedParams{
		UserID: int32(req.UserID),
		IsVerified: sql.NullBool{
			Bool:  req.IsVerfied,
			Valid: true,
		},
	})

	if err != nil {
		return nil, userrepositoryerrors.ErrUpdateUserVerificationCode
	}

	return r.mapper.ToUserRecord(res), nil
}

func (r *userCommandRepository) UpdateUserPassword(ctx context.Context, req *requests.UpdateUserPasswordRequest) (*record.UserRecord, error) {
	res, err := r.db.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:   int32(req.UserID),
		Password: req.Password,
	})

	if err != nil {
		return nil, userrepositoryerrors.ErrUpdateUserPassword
	}

	return r.mapper.ToUserRecord(res), nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.TrashUser(ctx, int32(user_id))

	if err != nil {
		return nil, userrepositoryerrors.ErrTrashedUser
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.RestoreUser(ctx, int32(user_id))

	if err != nil {
		return nil, userrepositoryerrors.ErrRestoreUser
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(ctx, int32(user_id))

	if err != nil {
		return false, userrepositoryerrors.ErrDeleteUserPermanent
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)

	if err != nil {
		return false, userrepositoryerrors.ErrRestoreAllUsers
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentUsers(ctx)

	if err != nil {
		return false, userrepositoryerrors.ErrDeleteAllUsers
	}
	return true, nil
}
