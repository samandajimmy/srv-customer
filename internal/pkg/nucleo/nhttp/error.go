package nhttp

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"

// Standard error response codes
const (
	BadRequestErrorCode          = "400"
	UnauthorizedErrorCode        = "401"
	ForbiddenErrorCode           = "403"
	NotFoundErrorCode            = "404"
	MethodNotAllowedErrorCode    = "405"
	UnprocessableEntityErrorCode = "422"
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
		HTTPStatusRespKey: 400,
	},
}

var UnprocessableEntityError = &ncore.Response{
	Success: false,
	Code:    UnprocessableEntityErrorCode,
	Message: "Unprocessable Entity",
	Metadata: map[string]interface{}{
		HTTPStatusRespKey: 422,
	},
}

var UnauthorizedError = &ncore.Response{
	Success: false,
	Code:    UnauthorizedErrorCode,
	Message: "Unauthorized",
	Metadata: map[string]interface{}{
		HTTPStatusRespKey: 401,
	},
}

var ForbiddenError = &ncore.Response{
	Success: false,
	Code:    ForbiddenErrorCode,
	Message: "Forbidden",
	Metadata: map[string]interface{}{
		HTTPStatusRespKey: 403,
	},
}

var NotFoundError = &ncore.Response{
	Success: false,
	Code:    NotFoundErrorCode,
	Message: "Not Found",
	Metadata: map[string]interface{}{
		HTTPStatusRespKey: 404,
	},
}

var MethodNotAllowedError = &ncore.Response{
	Success: false,
	Code:    MethodNotAllowedErrorCode,
	Message: "Method Not Allowed",
	Metadata: map[string]interface{}{
		HTTPStatusRespKey: 405,
	},
}
