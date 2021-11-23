package nhttp

import (
	"encoding/json"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

const (
	ContentTypeJSON = "application/json; charset=utf-8"
)

type JSONContentWriter struct {
	Debug bool
}

func (jw JSONContentWriter) Write(w http.ResponseWriter, httpStatus int, body interface{}) int {
	// Add content type
	w.Header().Add(ContentTypeHeader, ContentTypeJSON)
	// Write http status
	w.WriteHeader(httpStatus)
	// Send JSON response
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Errorf("failed to write response to json ( payload = %+v )", body)
	}
	// Return httpStatus
	return httpStatus
}

func (jw JSONContentWriter) WriteError(w http.ResponseWriter, err error) int {
	apiErr, ok := err.(*ncore.Response)
	if !ok {
		// If assert type fail, create wrap error to an internal error
		apiErr = ncore.InternalError.Wrap(err)
	}

	// Get http status
	httpStatus, ok := nval.ParseInt(apiErr.Metadata[HttpStatusRespKey])
	if !ok {
		httpStatus = http.StatusInternalServerError
	}

	// Get metadata of error
	metadata, _ := apiErr.Metadata[MetadataKey].(map[string]interface{})

	// Create response
	resp := Response{
		Success: false,
		Code:    apiErr.Code,
		Message: apiErr.Message,
		Data:    nil,
	}

	// If debug mode, then create error debug data
	if jw.Debug {
		// Get response message from source if exist
		respMessage := ""
		if apiErr.SourceError != nil {
			respMessage = apiErr.SourceError.Error()
		} else {
			respMessage = apiErr.Message
		}

		// Add error tracing metadata to data
		resp.Data = errorDataResponse{ErrorDebug: &errorDebug{
			Message:  respMessage,
			Traces:   apiErr.Traces,
			Metadata: metadata,
		}}
	}

	// Send error json
	return jw.Write(w, httpStatus, resp)
}
