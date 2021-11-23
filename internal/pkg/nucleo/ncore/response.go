package ncore

import (
	"fmt"
	"runtime"
)

const (
	SuccessResponseCode       = "OK"
	InternalErrorResponseCode = "ERROR"
)

type Response struct {
	Success     bool                   `yaml:"success" json:"success"`
	Code        string                 `yaml:"code" json:"code"`
	Message     string                 `yaml:"message" json:"message"`
	Metadata    map[string]interface{} `yaml:"metadata" json:"metadata"`
	SourceError error                  `yaml:"-" json:"-"`
	Traces      []string               `yaml:"-" json:"-"`
}

/// Error implements builtin error interface
func (r *Response) Error() string {
	if r.Success {
		return ""
	}
	errMsg := fmt.Sprintf("[%s] %s", r.Code, r.Message)
	if r.SourceError != nil {
		errMsg += fmt.Sprintf(" => %s", r.SourceError.Error())
	}
	return errMsg
}

// Unwrap / implements builtin error unwrapping
func (r *Response) Unwrap() error {
	return r.SourceError
}

func (r *Response) AddMetadata(k string, v interface{}) *Response {
	r.Metadata[k] = v
	return r
}

func NewError(str string) *Response {
	return &Response{
		Success: false,
		Code:    InternalErrorResponseCode,
		Message: str,
		Traces:  []string{Trace(1)},
	}
}

func (r Response) Wrap(sourceErr error, args ...interface{}) *Response {
	// Get arguments
	var metadataArg map[string]interface{}
	var skipTrace = 0
	switch len(args) {
	case 1:
		switch v := args[0].(type) {
		case int:
			skipTrace = v
		case map[string]interface{}:
			metadataArg = v
		}
	case 2:
		// args = metadata, skipTrace
		metadataArg, _ = args[0].(map[string]interface{})
		tmp, ok := args[1].(int)
		if ok {
			skipTrace = tmp
		}
	}

	// Reset skip trace to 1
	if skipTrace == 0 {
		skipTrace = 1
	}

	// Trace runtime caller
	trace := Trace(skipTrace)
	var traces []string
	if r.Traces == nil {
		// If traces is empty, then init traces
		traces = []string{trace}
	} else {
		// Else, push trace to stack traces
		traces = append(r.Traces, trace)
	}

	// Get response original metadata
	metadata := make(map[string]interface{})
	if r.Metadata != nil {
		for k, v := range r.Metadata {
			metadata[k] = v
		}
	}

	// Merge metadata from arguments
	for k, v := range metadataArg {
		metadata[k] = v
	}

	// Copy responses
	return &Response{
		Success:     false,
		Code:        r.Code,
		Message:     r.Message,
		Metadata:    metadata,
		SourceError: sourceErr,
		Traces:      traces,
	}
}

type ErrorOptions struct {
	SourceError error
	Metadata    map[string]interface{}
}

func Trace(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "<?>:<?>"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func TraceError(err error) error {
	if err == nil {
		return nil
	}

	// Check if err is Response
	respErr, ok := err.(*Response)
	if !ok {
		respErr = &Response{
			Success:     false,
			Code:        InternalError.Code,
			Message:     InternalError.Message,
			Metadata:    nil,
			SourceError: err,
			Traces:      []string{Trace(1)},
		}
		return respErr
	}

	// Append trace
	respErr.Traces = append(respErr.Traces, Trace(1))
	return respErr
}
