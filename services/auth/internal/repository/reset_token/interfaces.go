package resettokenrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
)

type ResetTokenQueryRepository interface {
	FindByResetToken(ctx context.Context, token string) (*record.ResetTokenRecord, error)
}

type ResetTokenCommandRepository interface {
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}
