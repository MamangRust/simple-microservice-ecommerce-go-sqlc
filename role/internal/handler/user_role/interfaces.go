package userrolehandler

import pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"

type UserRoleHandleGrpc interface {
	pbuserrole.UserRoleServiceServer
}
