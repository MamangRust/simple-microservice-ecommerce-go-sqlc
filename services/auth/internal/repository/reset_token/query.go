package resettokenrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
)

type resetTokenQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.ResetTokenRecordMapping
}

func NewResetTokenQueryRepository(db *db.Queries, mapper recordmapper.ResetTokenRecordMapping) *resetTokenQueryRepository {
	return &resetTokenQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *resetTokenQueryRepository) FindByResetToken(ctx context.Context, code string) (*record.ResetTokenRecord, error) {
	res, err := r.db.GetResetToken(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToResetTokenRecord(res), nil
}
