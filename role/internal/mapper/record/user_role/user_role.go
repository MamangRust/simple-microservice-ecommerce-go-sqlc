package userrolerecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type userRoleRecordMapper struct {
}

func NewUserRoleRecordMapper() UserRoleRecordMapping {
	return &userRoleRecordMapper{}
}

func (u *userRoleRecordMapper) ToUserRoleRecord(userRole *db.UserRole) *record.UserRoleRecord {
	return &record.UserRoleRecord{
		UserRoleID: int32(userRole.UserRoleID),
		UserID:     int32(userRole.UserID),
		RoleID:     int32(userRole.RoleID),
		CreatedAt:  userRole.CreatedAt.Time,
		UpdatedAt:  userRole.UpdatedAt.Time,
	}
}
