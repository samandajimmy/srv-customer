package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var log = nlogger.Get()

type DataSourceMap struct {
	DBInternal *nsql.DB
	DBExternal *nsql.DB
}

func (a *DataSourceMap) Init(config contract.DataSourcesConfig) error {
	// Skip if not initialized
	if a.DBInternal == nil {
		log.Debug("Skipping db init")
		return nil
	}

	// Init using prefix key on env
	err := a.DBInternal.Init(config.DBInternal, "")
	if err != nil {
		return ncore.TraceError(err)
	}

	err = a.DBExternal.Init(config.DBExternal, "EXTERNAL")
	if err != nil {
		return ncore.TraceError(err)
	}

	return nil
}

func NewDataSourceMap() DataSourceMap {
	return DataSourceMap{
		DBInternal: new(nsql.DB),
		DBExternal: new(nsql.DB),
	}
}
