package nsql

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

func IsUpdated(result sql.Result) bool {
	// Get affected rows
	count, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Errorf("nsql: unable to get affected rows"))
	}

	return count > 0
}

func NullTimeEpoch(t pq.NullTime) int64 {
	if !t.Valid {
		return 0
	}

	return t.Time.Unix()
}
