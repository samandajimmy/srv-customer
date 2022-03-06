package customer

import (
	"context"
	"encoding/base64"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func NewRepositoryExternal(config *DatabaseConfig) (*RepositoryExternal, error) {
	// Password Decode
	password, _ := base64.StdEncoding.DecodeString(config.DatabasePass)

	// Init db
	db, err := nsql.NewDatabase(nsql.Config{
		Driver:          config.DatabaseDriver,
		Host:            config.DatabaseHost,
		Port:            config.DatabasePort,
		Username:        config.DatabaseUser,
		Password:        string(password),
		Database:        config.DatabaseName,
		MaxIdleConn:     config.DatabaseMaxIdleConn,
		MaxOpenConn:     config.DatabaseMaxOpenConn,
		MaxConnLifetime: config.DatabaseMaxConnLifetime,
	})
	if err != nil {
		return nil, errx.Trace(err)
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
	r.InitializeStatement(ctx)

	// Get connection
	conn, err := r.db.GetConnection(ctx)
	if err != nil {
		log.Error("failed to retrieve connection to db", nlogger.Error(err))
		panic(errx.Trace(err))
	}

	return &RepositoryContext{
		ctx:  ctx,
		conn: conn,
		stmt: r.stmt,
		log:  nlogger.Get().NewChild(nlogger.Context(ctx)),
	}
}

func (r *RepositoryExternal) InitializeStatement(ctx context.Context) {
	// If db is not connected, then initialize connection
	isConnected, _ := r.db.IsConnected(ctx)
	if !isConnected {
		log.Debugf("initialize connection to database external...")
		err := r.db.InitContext(ctx)
		if err != nil {
			log.Error("failed to initiate connection to db", nlogger.Error(err))
			panic(errx.Trace(err))
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
}
