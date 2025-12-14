package rolerecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type roleQueryRecordMapper struct{}

func NewRoleQueryRecordMapper() RoleQueryRecordMapper {
	return &roleQueryRecordMapper{}
}

func (s *roleQueryRecordMapper) ToRoleRecord(role *db.Role) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleQueryRecordMapper) ToRolesRecord(roles []*db.Role) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecord(role))
	}

	return result
}

func (s *roleQueryRecordMapper) ToRoleRecordAll(role *db.GetRolesRow) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleQueryRecordMapper) ToRolesRecordAll(roles []*db.GetRolesRow) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecordAll(role))
	}

	return result
}

func (s *roleQueryRecordMapper) ToRoleRecordActive(role *db.GetActiveRolesRow) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleQueryRecordMapper) ToRolesRecordActive(roles []*db.GetActiveRolesRow) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecordActive(role))
	}

	return result
}

func (s *roleQueryRecordMapper) ToRoleRecordTrashed(role *db.GetTrashedRolesRow) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleQueryRecordMapper) ToRolesRecordTrashed(roles []*db.GetTrashedRolesRow) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecordTrashed(role))
	}

	return result
}
