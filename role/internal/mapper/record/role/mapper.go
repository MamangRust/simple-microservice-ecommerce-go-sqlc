package rolerecordmapper

type RoleRecordMapper interface {
	RoleQueryRecordMapper() RoleQueryRecordMapper
	RoleCommandRecordMapper() RoleCommandRecordMapper
}

type roleRecordMapper struct {
	roleQuery   RoleQueryRecordMapper
	roleCommand RoleCommandRecordMapper
}

func (r *roleRecordMapper) RoleQueryRecordMapper() RoleQueryRecordMapper {
	return r.roleQuery
}

func (r *roleRecordMapper) RoleCommandRecordMapper() RoleCommandRecordMapper {
	return r.roleCommand
}

func NewRoleRecordMapper() RoleRecordMapper {
	return &roleRecordMapper{
		roleQuery:   NewRoleQueryRecordMapper(),
		roleCommand: NewRoleCommandRecordMapper(),
	}
}
