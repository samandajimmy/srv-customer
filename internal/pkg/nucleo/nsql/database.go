package nsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

func NewDatabase(config Config) (*Database, error) {
	// Set default connection values
	config.normalizeValue()

	// Generate DSN
	dsn, err := config.getDSN()
	if err != nil {
		return nil, err
	}

	// Set config
	db := Database{
		config: &config,
		dsn:    dsn,
	}
	return &db, nil
}

type Database struct {
	config *Config
	dsn    string
	conn   *sqlx.DB
}

// Prepare prepare sql statements or exit app if fails or error
func (s *Database) Prepare(query string) *sqlx.Stmt {
	stmt, err := s.conn.Preparex(query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

// PrepareFmt prepare sql statements from string format or exit app if fails or error
func (s *Database) PrepareFmt(queryFmt string, args ...interface{}) *sqlx.Stmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.Prepare(query)
}

// PrepareNamedFmt prepare sql statements from string format with named bindvars or exit app if fails or error
func (s *Database) PrepareNamedFmt(queryFmt string, args ...interface{}) *sqlx.NamedStmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.PrepareNamed(query)
}

// PrepareNamed prepare sql statements with named bindvars or exit app if fails or error
func (s *Database) PrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := s.conn.PrepareNamed(query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

// PrepareTemplate prepare sql statements from a string template format or exit app if fails or error
func (s *Database) PrepareTemplate(q string, values map[string]string) *sqlx.Stmt {
	for a, v := range values {
		q = strings.ReplaceAll(q, ":"+a, v)
	}
	return s.Prepare(q)
}

// ReleaseTx clean db transaction by commit if no error, or rollback if an error occurred
func (s *Database) ReleaseTx(tx *sqlx.Tx, err *error) {
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

func (s *Database) Init() error {
	return s.InitContext(context.Background())
}

func (s *Database) InitContext(ctx context.Context) error {
	// Create connection
	conn, err := sqlx.ConnectContext(ctx, s.config.Driver, s.dsn)
	if err != nil {
		return err
	}

	// Set connection settings
	conn.SetConnMaxLifetime(time.Duration(*s.config.MaxConnLifetime) * time.Second)
	conn.SetMaxOpenConns(*s.config.MaxOpenConn)
	conn.SetMaxIdleConns(*s.config.MaxIdleConn)

	// Set connection
	s.conn = conn

	return nil
}

func (s *Database) IsConnected(ctx context.Context) (bool, error) {
	if s.conn == nil {
		return false, nil
	}

	// Ping to database
	err := s.conn.PingContext(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Database) GetConnection(ctx context.Context) (*sqlx.Conn, error) {
	return s.conn.Connx(ctx)
}

func (s *Database) WithContext(ctx context.Context) *DatabaseContext {
	return &DatabaseContext{
		conn: s.conn,
		ctx:  ctx,
	}
}

func (s *Database) GetStats() sql.DBStats {
	return s.conn.Stats()
}
