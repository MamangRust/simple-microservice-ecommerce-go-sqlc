package refreshtokenrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
	refreshtokenrepositoryerror "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/refresh_token/repository"
)

type refreshTokenQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.RefreshTokenRecordMapping
}

func NewRefreshTokenQueryRepository(db *db.Queries, mapping recordmapper.RefreshTokenRecordMapping) *refreshTokenQueryRepository {
	return &refreshTokenQueryRepository{}
}

func (r *refreshTokenQueryRepository) FindByToken(ctx context.Context, token string) (*record.RefreshTokenRecord, error) {
	res, err := r.db.FindRefreshTokenByToken(ctx, token)

	if err != nil {
		return nil, refreshtokenrepositoryerror.ErrTokenNotFound
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenQueryRepository) FindByUserId(ctx context.Context, user_id int) (*record.RefreshTokenRecord, error) {
	res, err := r.db.FindRefreshTokenByUserId(ctx, int32(user_id))

	if err != nil {
		return nil, refreshtokenrepositoryerror.ErrFindByUserID
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}
