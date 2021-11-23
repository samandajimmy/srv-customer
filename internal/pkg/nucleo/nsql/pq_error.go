package nsql

import "github.com/lib/pq"

type ErrorMetadata struct {
	Constraint string
	Message    string
}

func GetPostgresError(err error) (ErrorCode, *ErrorMetadata) {
	// Cast error
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return UnknownError, nil
	}

	switch pqErr.Code {
	case "23505":
		return UniqueError, &ErrorMetadata{Constraint: pqErr.Constraint, Message: pqErr.Detail}
	case "23503":
		return FKViolationError, &ErrorMetadata{Constraint: pqErr.Constraint, Message: pqErr.Detail}
	default:
		return UnhandledError, nil
	}
}
