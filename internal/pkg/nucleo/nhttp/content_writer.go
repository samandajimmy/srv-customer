package nhttp

import "net/http"

type ContentWriter interface {
	Write(w http.ResponseWriter, httpStatus int, body interface{})
	WriteError(w http.ResponseWriter, err error) int
}
