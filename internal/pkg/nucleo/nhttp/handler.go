package nhttp

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"time"
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

	// Call handler function
	result, err := h.fn(&rx)

	// Determine result
	var httpStatus int

	// If an error occurred, then write Error
	if err != nil {
		// If an error occurred, then write error response
		httpStatus = h.contentWriter.WriteError(w, err)
		rx.SetContextValue(HttpStatusRespKey, httpStatus)
		return
	}

	// If result is nil, then return No Content
	if result == nil {
		httpStatus = http.StatusNoContent
		w.WriteHeader(httpStatus)
		rx.SetContextValue(HttpStatusRespKey, httpStatus)
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
		result.Code = ncore.Success.Code
	}

	// If Code is not set, then set to success message
	if result.Message == "" {
		result.Message = ncore.Success.Message
	}

	// Set status http by result code
	var statusHttp int
	switch result.Code {
	case "200":
		statusHttp = http.StatusOK
	case "400":
		statusHttp = http.StatusBadRequest
	case "422":
		statusHttp = http.StatusUnprocessableEntity
	case "503":
		statusHttp = http.StatusServiceUnavailable
	default:
		statusHttp = http.StatusOK
	}

	if result.responseFlag == ViewRequest {
		// Write response html
		httpStatus = h.contentWriter.WriteView(w, statusHttp, result.Data)
	} else {
		// Write standard response json
		httpStatus = h.contentWriter.Write(w, statusHttp, result)
	}

	rx.SetContextValue(HttpStatusRespKey, httpStatus)
}

func HandleErrorNotFound(_ *Request) (*Response, error) {
	return nil, NotFoundError
}

func HandleErrorMethodNotAllowed(_ *Request) (*Response, error) {
	return nil, MethodNotAllowedError
}

func NewAppStatusHandler(startedAt time.Time, version string, args ...ContentWriter) http.Handler {
	// Get content writer from args
	var cw ContentWriter
	if len(args) == 0 {
		// If content writer is not set, set to JSON Content Writer
		cw = new(JSONContentWriter)
	} else {
		cw = args[0]
	}

	// Create handler
	return NewHandler(func(r *Request) (*Response, error) {
		// Compose response
		resp := OK()
		resp.Data = map[string]string{
			"uptime":  time.Since(startedAt).String(),
			"version": version,
		}
		return resp, nil
	}).SetWriter(cw)
}
