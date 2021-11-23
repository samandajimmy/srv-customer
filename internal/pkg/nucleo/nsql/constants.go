package nsql

const (
	DriverMySQL      = "mysql"
	DriverPostgreSQL = "postgres"

	Null = `null`
)

type ErrorCode = int8

const (
	UnknownError = ErrorCode(iota)
	UnhandledError
	UniqueError
	FKViolationError
)
