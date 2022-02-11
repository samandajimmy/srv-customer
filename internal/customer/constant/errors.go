package constant

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"

// Common Error

var ResourceNotFoundError = ncore.NewTraceableError("E_COMM_1", "Resource not found")
var StaleResourceError = ncore.NewTraceableError("E_COMM_2", "Cannot update stale resource")
