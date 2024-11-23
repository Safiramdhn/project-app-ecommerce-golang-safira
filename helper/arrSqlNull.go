package helper

import "database/sql"

func ConvertToNullInt64Slice(intSlice []int) []sql.NullInt64 {
	nullInt64Slice := make([]sql.NullInt64, len(intSlice))
	for i, val := range intSlice {
		nullInt64Slice[i] = sql.NullInt64{
			Int64: int64(val),
			Valid: true,
		}
	}
	return nullInt64Slice
}
