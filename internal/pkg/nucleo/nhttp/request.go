package nhttp

import (
	"context"
	"encoding/json"
	"net/http"
)

func NewRequest(r *http.Request) Request {
	return Request{
		Request:      r,
		hookMetadata: make(map[string]interface{}),
	}
}

type Request struct {
	*http.Request
	hookMetadata map[string]interface{}
}

func (r *Request) ParseJSONBody(dest interface{}) error {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(dest); err != nil {
		return err
	}
	return nil
}

/// SetContextValue set value to context.Context in http.Request
/// Value is accessible chain of http.Handler
func (r *Request) SetContextValue(k string, v interface{}) {
	ctx := context.WithValue(r.Context(), k, v)
	*r.Request = *r.WithContext(ctx)
}

/// GetContextValue get value to context.Context
/// Value is accessible chain of http.Handler
func (r *Request) GetContextValue(k string) interface{} {
	return r.Context().Value(k)
}

/// SetMetadata set value to Request metadata
/// Value is only accessible in one nhttp.Handler and before and after nhttp.HookFunc
func (r *Request) SetMetadata(k string, v interface{}) {
	r.hookMetadata[k] = v
}

/// GetMetadata get value to Request metadata
/// Value is only accessible in one nhttp.Handler and before and after nhttp.HookFunc
func (r *Request) GetMetadata(k string) interface{} {
	return r.hookMetadata[k]
}

/// End set value to Request context.Context to flag that this connection is ending and the next Handler will not
/// continue
func (r *Request) End(httpStatus int) {
	r.SetContextValue(HttpStatusRespKey, httpStatus)
}

func (r *Request) HasEnded() bool {
	v := r.GetContextValue(HttpStatusRespKey)
	return v != nil
}

func (r *Request) GetClientIP() string {
	metadata, ok := r.Context().Value(RequestMetadataKey).(RequestMetadata)
	if !ok {
		return "N/A"
	}

	return metadata.ClientIP
}

func ParseJSONBody(dest interface{}, r *http.Request) error {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(dest); err != nil {
		return BadRequestError.Wrap(err)
	}
	return nil
}

func WriteJSONError(w http.ResponseWriter, httpStatus int, body interface{}) {
	// Add content type
	w.Header().Add(ContentTypeHeader, ContentTypeJSON)
	// Write http status
	w.WriteHeader(httpStatus)
	// Send JSON response
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Errorf("failed to write response to json ( payload = %+v )", body)
	}
}
