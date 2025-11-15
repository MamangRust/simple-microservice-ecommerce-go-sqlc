package resettokenrepository

import (
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
)

type ResetTokenRepository interface {
	ResetTokenQueryRepository
	ResetTokenCommandRepository
}

type resetTokenRepository struct {
	*resetTokenQueryRepository
	*resetTokenCommandRepository
}

func NewResetTokenRepository(db *db.Queries) ResetTokenRepository {
	mapper := recordmapper.NewResetTokenRecordMapper()

	return &resetTokenRepository{
		resetTokenQueryRepository:   NewResetTokenQueryRepository(db, mapper),
		resetTokenCommandRepository: NewResetTokenCommandRepository(db, mapper),
	}
}
