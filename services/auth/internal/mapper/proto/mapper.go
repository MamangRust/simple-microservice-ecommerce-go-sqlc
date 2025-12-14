package protomapper

type ProtoMapper struct {
	AuthProtoMapper AuthProtoMapper
	UserProtoMapper UserProtoMapper
}

func NewProtoMapper() *ProtoMapper {
	return &ProtoMapper{
		AuthProtoMapper: NewAuthProtoMapper(),
		UserProtoMapper: NewUserProtoMapper(),
	}
}
