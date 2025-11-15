package recordmapper

type RecordMapper struct {
	RefreshTokenMapper RefreshTokenRecordMapping
	ResetTokenMapper   ResetTokenRecordMapping
}

func NewRecordMapper() *RecordMapper {
	return &RecordMapper{
		RefreshTokenMapper: NewRefreshTokenRecordMapper(),
		ResetTokenMapper:   NewResetTokenRecordMapper(),
	}
}
