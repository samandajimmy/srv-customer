package customer

import (
	"context"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func NewRepositoryExternal(config *DatabaseConfig) (*RepositoryExternal, error) {
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
		return nil, ncore.TraceError("failed to init database", err)
	}

	// Init repo
	r := RepositoryExternal{
		db: db,
	}

	return &r, nil
}

type RepositoryExternal struct {
	db   *nsql.Database
	stmt *statement.Statements
}

func (r *RepositoryExternal) WithContext(ctx context.Context) *RepositoryContext {
	// If db is not connected, then initialize connection
	isConnected, _ := r.db.IsConnected(ctx)
	if !isConnected {
		log.Debugf("initialize connection to database...")
		err := r.db.InitContext(ctx)
		if err != nil {
			log.Error("failed to initiate connection to db", nlogger.Error(err))
			panic(ncore.TraceError("failed to initiate connection to db", err))
		}

		// Empty statements to re-initialize
		r.stmt = nil
	}

	// Initialize statements
	if r.stmt == nil {
		log.Debugf("initialize statement...")
		dbc := r.db.WithContext(ctx)
		r.stmt = statement.NewExternal(dbc)
	}

	// Get connection
	conn, err := r.db.GetConnection(ctx)
	if err != nil {
		log.Error("failed to retrieve connection to db", nlogger.Error(err))
		panic(ncore.TraceError("failed to retrieve connection to db", err))
	}

	return &RepositoryContext{
		ctx:  ctx,
		conn: conn,
		stmt: r.stmt,
		log:  nlogger.Get().NewChild(nlogger.Context(ctx)),
	}
}
