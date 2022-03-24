package customer

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func NewRepository(config *DatabaseConfig) (*Repository, error) {
	// Init db
	db, err := nsql.NewDatabase(nsql.Config{
		Driver:          config.DatabaseDriver,
		Host:            config.DatabaseHost,
		Port:            config.DatabasePort,
		Username:        config.DatabaseUser,
		Password:        config.DatabasePass,
		Database:        config.DatabaseName,
		MaxIdleConn:     config.DatabaseMaxIdleConn,
		MaxOpenConn:     config.DatabaseMaxOpenConn,
		MaxConnLifetime: config.DatabaseMaxConnLifetime,
	})
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Init repo
	r := Repository{
		db: db,
	}

	return &r, nil
}

type Repository struct {
	db   *nsql.Database
	stmt *statement.Statements
}

func (r *Repository) WithContext(ctx context.Context) *RepositoryContext {
	r.InitializeStatement(ctx)

	// Get connection
	conn, err := r.db.GetConnection(ctx)
	if err != nil {
		log.Error("failed to retrieve connection to db", logOption.Error(err))
		panic(errx.Trace(err))
	}

	return &RepositoryContext{
		ctx:  ctx,
		conn: conn,
		stmt: r.stmt,
		log:  nlogger.NewChild(logOption.WithNamespace("repository"), logOption.Context(ctx)),
	}
}

func (r *Repository) InitializeStatement(ctx context.Context) {
	// If db is not connected, then initialize connection
	isConnected, _ := r.db.IsConnected(ctx)
	if !isConnected {
		log.Debugf("initialize connection to database...")
		err := r.db.InitContext(ctx)
		if err != nil {
			log.Error("failed to initiate connection to db", logOption.Error(err))
			panic(errx.Trace(err))
		}

		// Empty statements to re-initialize
		r.stmt = nil
	}

	// Initialize statements
	if r.stmt == nil {
		log.Debugf("initialize statement...")
		dbc := r.db.WithContext(ctx)
		r.stmt = statement.New(dbc)
	}
}

type RepositoryContext struct {
	stmt *statement.Statements
	ctx  context.Context
	conn *sqlx.Conn
	log  nlogger.Logger
}

// ReleaseTx clean db transaction by commit if no error, or rollback if an error occurred
func (rc *RepositoryContext) ReleaseTx(tx *sqlx.Tx, err *error) {
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
