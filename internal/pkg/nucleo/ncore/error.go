package ncore

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	HttpStatusKey = "httpStatus"
)

var InternalErrorMetadata = map[string]interface{}{
	HttpStatusKey: 500,
}

// Trace returns where in file and line the function being called
func Trace(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "<?>:<?>"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func NewTraceableError(code string, message string) *TraceableError {
	return &TraceableError{
		Code:     code,
		Message:  message,
		Metadata: make(map[string]interface{}),
	}
}

// TraceableError print error meaningful message and stack trace for easier error tracing on certain condition such as
// panic
type TraceableError struct {
	Code        string                 `yaml:"code" json:"code"`
	Message     string                 `yaml:"message" json:"message"`
	Metadata    map[string]interface{} `yaml:"metadata" json:"metadata"`
	SourceError error                  `yaml:"-" json:"-"`
	Traces      []string               `yaml:"-" json:"-"`
}

func (t *TraceableError) GetTraces() []string {
	return t.Traces
}

func (t *TraceableError) AddMetadata(k string, v interface{}) *TraceableError {
	if t.Metadata == nil {
		t.Metadata = make(map[string]interface{})
	}
	t.Metadata[k] = v

	return t
}

func (t *TraceableError) Wrap(source error) *TraceableError {
	wErr := t.Trace()
	t.SourceError = source
	return wErr
}

func (t *TraceableError) Error() string {
	sourceErr := ""
	if t.SourceError != nil {
		sourceErr = t.SourceError.Error()
	}
	return fmt.Sprintf("%s: %s.\n  CausedBy => %s\n  Traces => %s", t.Code, t.Message,
		sourceErr, strings.Join(t.Traces, "\n            "))
}

// Unwrap implements builtin error unwrapping
func (t *TraceableError) Unwrap() error {
	return t.SourceError
}

func (t *TraceableError) Trace() *TraceableError {
	if t.Traces == nil {
		t.Traces = []string{Trace(2)}
	}

	traces := append(t.Traces, Trace(1))

	// Copy error and add trace
	return &TraceableError{
		Code:        t.Code,
		Message:     t.Message,
		Metadata:    t.Metadata,
		SourceError: t.SourceError,
		Traces:      traces,
	}
}

func TraceError(message string, err error) *TraceableError {
	// If error is nil, then ignore
	if err == nil {
		return nil
	}

	// Init empty traceable error
	var tErr *TraceableError

	// If is a Response, then convert to a traceable error
	switch vErr := err.(type) {
	case *TraceableError:
		// Copy error
		tErr = &TraceableError{
			Code:        vErr.Code,
			Message:     vErr.Message,
			Metadata:    vErr.Metadata,
			SourceError: vErr.SourceError,
			Traces:      vErr.Traces,
		}
	default:
		tErr = &TraceableError{
			Code:        "internal",
			Message:     message,
			SourceError: err,
			Traces:      []string{Trace(2)},
			Metadata:    InternalErrorMetadata,
		}
	}

	// Push traces
	tErr.Traces = append(tErr.Traces, Trace(1))

	return tErr
}

func (t *TraceableError) MessageOnly() string {
	return fmt.Sprintf("%s: %s", t.Code, t.Message)
}
