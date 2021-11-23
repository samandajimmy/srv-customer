package nhttp

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"

const (
	// Standard error response codes
	BadRequestErrorCode       = "400"
	UnauthorizedErrorCode     = "401"
	ForbiddenErrorCode        = "403"
	NotFoundErrorCode         = "404"
	MethodNotAllowedErrorCode = "405"
)

type errorDataResponse struct {
	ErrorDebug *errorDebug `json:"_error,omitempty"`
}

type errorDebug struct {
	Message  string      `json:"message,omitempty"`
	Traces   []string    `json:"traces,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

var BadRequestError = &ncore.Response{
	Success: false,
	Code:    BadRequestErrorCode,
	Message: "Bad Request",
	Metadata: map[string]interface{}{
		HttpStatusRespKey: 400,
	},
}

var UnauthorizedError = &ncore.Response{
	Success: false,
	Code:    UnauthorizedErrorCode,
	Message: "Unauthorized",
	Metadata: map[string]interface{}{
		HttpStatusRespKey: 401,
	},
}

var ForbiddenError = &ncore.Response{
	Success: false,
	Code:    ForbiddenErrorCode,
	Message: "Forbidden",
	Metadata: map[string]interface{}{
		HttpStatusRespKey: 403,
	},
}

var NotFoundError = &ncore.Response{
	Success: false,
	Code:    NotFoundErrorCode,
	Message: "Not Found",
	Metadata: map[string]interface{}{
		HttpStatusRespKey: 404,
	},
}

var MethodNotAllowedError = &ncore.Response{
	Success: false,
	Code:    MethodNotAllowedErrorCode,
	Message: "Method Not Allowed",
	Metadata: map[string]interface{}{
		HttpStatusRespKey: 405,
	},
}
