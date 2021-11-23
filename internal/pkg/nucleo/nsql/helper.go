package nsql

import "database/sql"

func IsUpdated(result sql.Result) error {
	// Get rows affected
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If not affected, then throw error
	if count == 0 {
		return RowNotUpdatedError
	}

	return nil
}
