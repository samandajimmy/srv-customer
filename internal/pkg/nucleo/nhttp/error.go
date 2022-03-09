package nhttp

import (
	"github.com/nbs-go/errx"
	"net/http"
)

const (
	HttpStatusMetadata = "httpStatus"
	MessageMetadata    = "message"
)

func WithStatus(status int) errx.SetOptionFn {
	return errx.AddMetadata(HttpStatusMetadata, status)
}

var b = errx.NewBuilder(pkgNamespace)

// Standard Errors

var InternalError = b.NewError("500", "Internal Error",
	WithStatus(http.StatusInternalServerError))

var BadRequestError = b.NewError("400", "Bad Request",
	WithStatus(http.StatusBadRequest),
)

var UnauthorizedError = b.NewError("401", "Unauthorized",
	WithStatus(http.StatusUnauthorized),
)

var ForbiddenError = b.NewError("403", "Forbidden",
	WithStatus(http.StatusForbidden),
)

var NotFoundError = b.NewError("404", "Not Found",
	WithStatus(http.StatusNotFound),
)

var MethodNotAllowedError = b.NewError("405", "Method Not Allowed",
	WithStatus(http.StatusMethodNotAllowed),
)

// Authorization Errors

var EmptyAuthorizationError = b.NewError("E_AUTH_1", "Authorization value is empty",
	WithStatus(http.StatusBadRequest))

var MalformedTokenError = errx.NewError("E_AUTH_2", "Malformed token",
	WithStatus(http.StatusBadRequest))
