package resettokenrepository

import (
	"context"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
	resettokenrepositoryerror "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/reset_token/repository"
)

type resetTokenCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.ResetTokenRecordMapping
}

func NewResetTokenCommandRepository(db *db.Queries, mapper recordmapper.ResetTokenRecordMapping) *resetTokenCommandRepository {
	return &resetTokenCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateResetToken inserts a new reset token into the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing user ID and token info
//
// Returns:
//   - The created ResetTokenRecord, or an error if the operation fails.
func (r *resetTokenCommandRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	res, err := r.db.CreateResetToken(ctx, db.CreateResetTokenParams{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	})
	if err != nil {
		return nil, resettokenrepositoryerror.ErrCreateResetToken
	}
	return r.mapper.ToResetTokenRecord(res), nil
}

// DeleteResetToken removes the reset token associated with the given user ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - userID: the user ID whose token should be deleted
//
// Returns:
//   - An error if the deletion fails.
func (r *resetTokenCommandRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	err := r.db.DeleteResetToken(ctx, int64(user_id))
	if err != nil {
		return resettokenrepositoryerror.ErrDeleteResetToken
	}
	return nil
}
