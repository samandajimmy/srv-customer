package handler

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"

type Middlewares struct {
	serviceMap *contract.ServiceMap
}

func NewMiddlewares(svc *contract.ServiceMap) *Middlewares {
	m := Middlewares{serviceMap: svc}
	return &m
}
