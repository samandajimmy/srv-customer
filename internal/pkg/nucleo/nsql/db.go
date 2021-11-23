package nsql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"time"
)

type DB struct {
	Conn *sqlx.DB
}

// Prepare prepare sql statements or exit app if fails or error
func (s *DB) Prepare(query string) *sqlx.Stmt {
	stmt, err := s.Conn.Preparex(query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

// PrepareFmt prepare sql statements from string format or exit app if fails or error
func (s *DB) PrepareFmt(queryFmt string, args ...interface{}) *sqlx.Stmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.Prepare(query)
}

// PrepareNamedFmt prepare sql statements from string format with named bindvars or exit app if fails or error
func (s *DB) PrepareNamedFmt(queryFmt string, args ...interface{}) *sqlx.NamedStmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.PrepareNamed(query)
}

// PrepareNamed prepare sql statements with named bindvars or exit app if fails or error
func (s *DB) PrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := s.Conn.PrepareNamed(query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

// ReleaseTx clean db transaction by commit if no error, or rollback if an error occurred
func (s *DB) ReleaseTx(tx *sqlx.Tx, err *error) {
	if *err != nil {
		// If an error occurred, rollback transaction
		errRollback := tx.Rollback()
		if errRollback != nil {
			panic(fmt.Errorf("failed to rollback database transaction.\n  > %w", errRollback))
		}
		return
	}

	// Else, commit transaction
	errCommit := tx.Commit()
	if errCommit != nil {
		panic(fmt.Errorf("failed to commit database transaction\n  > %w", errCommit))
	}
}

func (s *DB) Init(config Config) error {
	// Load configuration from env
	config.loadFromEnv()

	// Generate DSN
	dsn, err := config.getDSN()
	if err != nil {
		return ncore.TraceError(err)
	}

	// Create connection
	conn, err := sqlx.Connect(config.Driver, dsn)
	if err != nil {
		return ncore.TraceError(err)
	}

	// Set connection settings
	conn.SetConnMaxLifetime(time.Duration(*config.MaxConnLifetime) * time.Second)
	conn.SetMaxOpenConns(*config.MaxOpenConn)
	conn.SetMaxIdleConns(*config.MaxIdleConn)

	// Set connection
	s.Conn = conn

	return nil
}
