package recordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
)

type userQueryRecordMapper struct{}

func NewUserQueryRecordMapper() UserQueryRecordMapper {
	return &userQueryRecordMapper{}
}

func (s *userQueryRecordMapper) ToUserRecord(user *db.User) *record.UserRecord {
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

func (s *userQueryRecordMapper) ToUserRecordPagination(user *db.GetUsersRow) *record.UserRecord {
	var deletedAt *string

	if user.DeletedAt.Valid {
		formatedDeletedAt := user.DeletedAt.Time.Format("2006-01-02")

		deletedAt = &formatedDeletedAt
	}

	return &record.UserRecord{
		ID:         int(user.UserID),
		FirstName:  user.Firstname,
		LastName:   user.Lastname,
		Email:      user.Email,
		Password:   user.Password,
		IsVerified: user.IsVerified.Valid,
		CreatedAt:  user.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:  deletedAt,
	}
}

func (s *userQueryRecordMapper) ToUsersRecordPagination(users []*db.GetUsersRow) []*record.UserRecord {
	var userRecords []*record.UserRecord

	for _, user := range users {
		userRecords = append(userRecords, s.ToUserRecordPagination(user))
	}

	return userRecords
}

func (s *userQueryRecordMapper) ToUserRecordActivePagination(user *db.GetUsersActiveRow) *record.UserRecord {
	var deletedAt *string

	if user.DeletedAt.Valid {
		formatedDeletedAt := user.DeletedAt.Time.Format("2006-01-02")

		deletedAt = &formatedDeletedAt
	}

	return &record.UserRecord{
		ID:         int(user.UserID),
		FirstName:  user.Firstname,
		LastName:   user.Lastname,
		Email:      user.Email,
		Password:   user.Password,
		IsVerified: user.IsVerified.Valid,
		CreatedAt:  user.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:  deletedAt,
	}
}

func (s *userQueryRecordMapper) ToUsersRecordActivePagination(users []*db.GetUsersActiveRow) []*record.UserRecord {
	var userRecords []*record.UserRecord

	for _, user := range users {
		userRecords = append(userRecords, s.ToUserRecordActivePagination(user))
	}

	return userRecords
}

func (s *userQueryRecordMapper) ToUserRecordTrashedPagination(user *db.GetUserTrashedRow) *record.UserRecord {
	var deletedAt *string

	if user.DeletedAt.Valid {
		formatedDeletedAt := user.DeletedAt.Time.Format("2006-01-02")

		deletedAt = &formatedDeletedAt
	}

	return &record.UserRecord{
		ID:        int(user.UserID),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}

func (s *userQueryRecordMapper) ToUsersRecordTrashedPagination(users []*db.GetUserTrashedRow) []*record.UserRecord {
	var userRecords []*record.UserRecord

	for _, user := range users {
		userRecords = append(userRecords, s.ToUserRecordTrashedPagination(user))
	}

	return userRecords
}
