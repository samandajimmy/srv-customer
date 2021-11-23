package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type Initializer interface {
	ncore.InitializeChecker
	Init(dataSources DataSourceMap, repositories contract.RepositoryMap) error
}
