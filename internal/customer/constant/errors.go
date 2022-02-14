package constant

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"

var InternalError = ncore.NewTraceableError("500", "Internal Error")
var BadRequestError = ncore.NewTraceableError("400", "Bad Request")
var ForbiddenError = ncore.NewTraceableError("403", "Forbidden")

// Common Error

var ResourceNotFoundError = ncore.NewTraceableError("E_COMM_1", "Resource not found")
var StaleResourceError = ncore.NewTraceableError("E_COMM_2", "Cannot update stale resource")
var DefaultError = ncore.NewTraceableError("E_COMM_3", "An error has occurred, please try again later")
