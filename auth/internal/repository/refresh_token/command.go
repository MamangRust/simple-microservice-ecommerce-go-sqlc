package refreshtokenrepository

import (
	"context"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
	refreshtokenrepositoryerror "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/refresh_token/repository"
)

type refreshTokenCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.RefreshTokenRecordMapping
}

func NewRefreshTokenCommandRepository(db *db.Queries, mapper recordmapper.RefreshTokenRecordMapping) *refreshTokenCommandRepository {
	return &refreshTokenCommandRepository{db: db, mapper: mapper}
}

// CreateRefreshToken inserts a new refresh token record into the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing token and user information
//
// Returns:
//   - The created RefreshTokenRecord, or an error if the operation fails.
func (r *refreshTokenCommandRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtokenrepositoryerror.ErrParseDate
	}

	res, err := r.db.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})

	if err != nil {
		return nil, refreshtokenrepositoryerror.ErrCreateRefreshToken
	}

	return r.mapper.ToRefreshTokenRecord(res), nil
}

// UpdateRefreshToken updates an existing refresh token record.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload with updated token data
//
// Returns:
//   - The updated RefreshTokenRecord, or an error if the operation fails.
func (r *refreshTokenCommandRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtokenrepositoryerror.ErrParseDate
	}

	res, err := r.db.UpdateRefreshTokenByUserId(ctx, db.UpdateRefreshTokenByUserIdParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})
	if err != nil {
		return nil, refreshtokenrepositoryerror.ErrUpdateRefreshToken
	}

	return r.mapper.ToRefreshTokenRecord(res), nil
}

// DeleteRefreshToken removes a refresh token by its token string.
//
// Parameters:
//   - ctx: the context for the database operation
//   - token: the refresh token string to delete
//
// Returns:
//   - An error if the deletion fails.
func (r *refreshTokenCommandRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.db.DeleteRefreshToken(ctx, token)

	if err != nil {
		return refreshtokenrepositoryerror.ErrDeleteRefreshToken
	}

	return nil
}

// DeleteRefreshTokenByUserId removes a refresh token by the associated user ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - user_id: the ID of the user whose token will be deleted
//
// Returns:
//   - An error if the deletion fails.
func (r *refreshTokenCommandRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	err := r.db.DeleteRefreshTokenByUserId(ctx, int32(user_id))

	if err != nil {
		return refreshtokenrepositoryerror.ErrDeleteByUserID
	}

	return nil
}
