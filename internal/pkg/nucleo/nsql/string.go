package nsql

import "database/sql"

func NewNullString(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  str != "",
	}
}

func NewNullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: v,
		Valid: v != 0,
	}
}
