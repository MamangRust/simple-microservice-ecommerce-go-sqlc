package repository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
)

type UserQueryRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)
	FindByEmail(ctx context.Context, email string) (*record.UserRecord, error)

	FindByEmailAndVerify(ctx context.Context, email string) (*record.UserRecord, error)
	FindByVerificationCode(ctx context.Context, code string) (*record.UserRecord, error)
}

type UserCommandRepository interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*record.UserRecord, error)

	UpdateUserIsVerified(ctx context.Context, req *requests.UpdateUserVerifiedRequest) (*record.UserRecord, error)
	UpdateUserPassword(ctx context.Context, req *requests.UpdateUserPasswordRequest) (*record.UserRecord, error)

	TrashedUser(ctx context.Context, user_id int) (*record.UserRecord, error)
	RestoreUser(ctx context.Context, user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}
