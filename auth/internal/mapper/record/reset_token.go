package recordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
)

type resetTokenRecordMapper struct {
}

func NewResetTokenRecordMapper() *resetTokenRecordMapper {
	return &resetTokenRecordMapper{}
}

func (r *resetTokenRecordMapper) ToResetTokenRecord(resetToken *db.ResetToken) *record.ResetTokenRecord {
	return &record.ResetTokenRecord{
		ID:        int64(resetToken.ID),
		UserID:    resetToken.UserID,
		Token:     resetToken.Token,
		ExpiredAt: resetToken.ExpiryDate.String(),
	}
}
