package nsql

import "errors"

const (
	DriverMySQL      = "mysql"
	DriverPostgreSQL = "postgres"

	Null = `null`

	SetVersion      = `"version" = "version" + 1`
	UpdateCondition = `"id" = :id AND "version" = :version`
)

type ErrorCode = int8

const (
	UnknownError = ErrorCode(iota)
	UnhandledError
	UniqueError
	FKViolationError
)

var ErrNoRowUpdated = errors.New("nsql: no row updated")
