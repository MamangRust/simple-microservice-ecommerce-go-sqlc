package refreshtokenrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
)

type RefreshTokenQueryRepository interface {
	FindByToken(ctx context.Context, token string) (*record.RefreshTokenRecord, error)
	FindByUserId(ctx context.Context, user_id int) (*record.RefreshTokenRecord, error)
}

type RefreshTokenCommandRepository interface {
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}
