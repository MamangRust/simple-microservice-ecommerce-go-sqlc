package roleresponsemapper

type RoleResponseMapper interface {
	RoleQueryResponseMapper() RoleQueryResponseMapper
	RoleCommandResponseMapper() RoleCommandResponseMapper
}

type roleResponseMapper struct {
	roleQueryMapper   RoleQueryResponseMapper
	roleCommandMapper RoleCommandResponseMapper
}

func (r *roleResponseMapper) RoleQueryResponseMapper() RoleQueryResponseMapper {
	return r.roleQueryMapper
}

func (r *roleResponseMapper) RoleCommandResponseMapper() RoleCommandResponseMapper {
	return r.roleCommandMapper
}

func NewRoleResponseMapper() RoleResponseMapper {
	return &roleResponseMapper{
		roleQueryMapper:   NewRoleQueryResponseMapper(),
		roleCommandMapper: NewRoleCommandResponseMapper(),
	}
}
