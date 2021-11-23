package ncore

var Success = Response{
	Success: true,
	Code:    SuccessResponseCode,
	Message: "OK",
}

var InternalError = &Response{
	Success: false,
	Code:    InternalErrorResponseCode,
	Message: "Internal Error",
}
