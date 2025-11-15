package rolerecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type RoleBaseRecordMapper interface {
	ToRoleRecord(role *db.Role) *record.RoleRecord
}

type RoleQueryRecordMapper interface {
	RoleBaseRecordMapper
	ToRolesRecord(roles []*db.Role) []*record.RoleRecord
	ToRoleRecordAll(role *db.GetRolesRow) *record.RoleRecord
	ToRolesRecordAll(roles []*db.GetRolesRow) []*record.RoleRecord
	ToRoleRecordActive(role *db.GetActiveRolesRow) *record.RoleRecord
	ToRolesRecordActive(roles []*db.GetActiveRolesRow) []*record.RoleRecord

	ToRoleRecordTrashed(role *db.GetTrashedRolesRow) *record.RoleRecord

	ToRolesRecordTrashed(roles []*db.GetTrashedRolesRow) []*record.RoleRecord
}

type RoleCommandRecordMapper interface {
	RoleBaseRecordMapper
}
