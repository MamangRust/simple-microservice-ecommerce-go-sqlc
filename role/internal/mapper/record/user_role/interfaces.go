package userrolerecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type UserRoleRecordMapping interface {
	ToUserRoleRecord(userRole *db.UserRole) *record.UserRoleRecord
}
