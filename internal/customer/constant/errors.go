package constant

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"

// Common Error

var ResourceNotFoundError = ncore.NewTraceableError("E_COMM_1", "Resource not found")
var VersioningError = ncore.NewTraceableError("E_COMM_2", "Invalid resource version to update")
