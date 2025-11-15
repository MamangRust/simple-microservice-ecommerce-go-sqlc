package rolerecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type roleCommandMapper struct {
}

func NewRoleCommandRecordMapper() RoleCommandRecordMapper {
	return &roleCommandMapper{}
}

func (s *roleCommandMapper) ToRoleRecord(role *db.Role) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}
