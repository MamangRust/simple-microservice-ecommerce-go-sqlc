package recordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
)

type userCommandRecordMapper struct{}

func NewUserCommandRecordMapper() UserCommandRecordMapper {
	return &userCommandRecordMapper{}
}

func (s *userCommandRecordMapper) ToUserRecord(user *db.User) *record.UserRecord {
	var deletedAt *string

	if user.DeletedAt.Valid {
		formatedDeletedAt := user.DeletedAt.Time.Format("2006-01-02")

		deletedAt = &formatedDeletedAt
	}

	if user.VerificationCode == "" {
		user.VerificationCode = "null"
	}

	return &record.UserRecord{
		ID:           int(user.UserID),
		FirstName:    user.Firstname,
		LastName:     user.Lastname,
		VerifiedCode: user.VerificationCode,
		IsVerified:   user.IsVerified.Valid,
		Email:        user.Email,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:    deletedAt,
	}
}
