package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"time"
)

func NewCommon(startTime time.Time, appVersion, appBuildHash string) *Common {
	h := Common{
		startTime: startTime,
		version:   appVersion,
		buildHash: appBuildHash,
	}
	return &h
}

type Common struct {
	startTime time.Time
	version   string
	buildHash string
}

func (c *Common) GetAPIStatus(_ *nhttp.Request) (*nhttp.Response, error) {
	res := nhttp.Success().
		SetData(map[string]string{
			"version":    c.version,
			"build_hash": c.buildHash,
			"uptime":     time.Since(c.startTime).String(),
		})
	return res, nil
}
