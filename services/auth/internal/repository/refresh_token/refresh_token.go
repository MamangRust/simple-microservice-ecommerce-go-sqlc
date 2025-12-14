package refreshtokenrepository

import (
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
)

type RefreshTokenRepository interface {
	RefreshTokenQueryRepository
	RefreshTokenCommandRepository
}

type refreshTokenRepository struct {
	RefreshTokenQueryRepository
	RefreshTokenCommandRepository
}

func NewRefreshTokenRepository(db *db.Queries) RefreshTokenRepository {
	mapper := recordmapper.NewRefreshTokenRecordMapper()

	return &refreshTokenRepository{
		RefreshTokenQueryRepository:   NewRefreshTokenQueryRepository(db, mapper),
		RefreshTokenCommandRepository: NewRefreshTokenCommandRepository(db, mapper),
	}
}
