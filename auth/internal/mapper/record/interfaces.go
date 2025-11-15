package recordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
)

type RefreshTokenRecordMapping interface {
	ToRefreshTokenRecord(refreshToken *db.RefreshToken) *record.RefreshTokenRecord
	ToRefreshTokensRecord(refreshTokens []*db.RefreshToken) []*record.RefreshTokenRecord
}

type ResetTokenRecordMapping interface {
	ToResetTokenRecord(resetToken *db.ResetToken) *record.ResetTokenRecord
}
