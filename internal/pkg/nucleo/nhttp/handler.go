package nhttp

import (
	"fmt"
	"net/http"
)

// HandlerFunc represents function that will be called
type HandlerFunc func(r *Request) (*Response, error)

// NewHandler initiate a new handler that implements http.Handler interface
func NewHandler(fn HandlerFunc) *Handler {
	h := Handler{fn: fn}
	return &h
}

/// Handler handles HTTP request and send response as JSON

type Handler struct {
	// Private
	contentWriter ContentWriter
	fn            HandlerFunc
}

func (h *Handler) SetWriter(cw ContentWriter) *Handler {
	h.contentWriter = cw
	return h
}

/// ServeHTTP implement http.Handler interface to write success or error response in JSON
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Prepare extended request
	rx := NewRequest(r)

	// Init recovery handler
	panicked := true
	defer func() {
		if result := recover(); result != nil || panicked {
			log.Errorf("Panic => %v", result)

			var err error
			if rErr, ok := result.(error); ok {
				err = rErr
			} else {
				err = InternalError.Wrap(fmt.Errorf("unknwon result received from panic: %v", result))
			}

			httpStatus := h.contentWriter.WriteError(w, err)
			rx.SetContextValue(HTTPStatusRespContextKey, httpStatus)
		}
	}()

	// Call handler function
	result, err := h.fn(&rx)
	panicked = false

	// If an error occurred, then write Error
	if err != nil {
		// If an error occurred, then write error response
		httpStatus := h.contentWriter.WriteError(w, err)
		rx.SetContextValue(HTTPStatusRespContextKey, httpStatus)
		return
	}

	// If result is nil, then return No Content
	if result == nil {
		httpStatus := http.StatusNoContent
		w.WriteHeader(httpStatus)
		rx.SetContextValue(HTTPStatusRespContextKey, httpStatus)
		return
	}

	// If response flag is continue, then return
	if result.responseFlag == ContinueRequest {
		return
	}

	// Set headers
	for k, v := range result.Header {
		w.Header().Set(k, v)
	}

	// Set standard response if not set
	if !result.Success {
		result.Success = true
	}

	// If Code is not set, then set to success code
	if result.Code == "" {
		result.Code = SuccessCode
	}

	// If Code is not set, then set to success message
	if result.Message == "" {
		result.Message = SuccessMessage
	}

	// Write standard response json
	httpStatus := http.StatusOK
	h.contentWriter.Write(w, httpStatus, result)
	rx.SetContextValue(HTTPStatusRespContextKey, httpStatus)
}

func HandleErrorNotFound(_ *Request) (*Response, error) {
	return nil, NotFoundError
}

func HandleErrorMethodNotAllowed(_ *Request) (*Response, error) {
	return nil, MethodNotAllowedError
}
