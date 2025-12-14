package recordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
)

type UserBaseRecordMapper interface {
	ToUserRecord(user *db.User) *record.UserRecord
}

type UserQueryRecordMapper interface {
	UserBaseRecordMapper

	ToUserRecordPagination(user *db.GetUsersRow) *record.UserRecord

	ToUsersRecordPagination(users []*db.GetUsersRow) []*record.UserRecord

	ToUserRecordActivePagination(user *db.GetUsersActiveRow) *record.UserRecord

	ToUsersRecordActivePagination(users []*db.GetUsersActiveRow) []*record.UserRecord

	ToUserRecordTrashedPagination(user *db.GetUserTrashedRow) *record.UserRecord

	ToUsersRecordTrashedPagination(users []*db.GetUserTrashedRow) []*record.UserRecord
}

type UserCommandRecordMapper interface {
	UserBaseRecordMapper
}
