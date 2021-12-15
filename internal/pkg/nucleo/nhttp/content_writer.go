package nhttp

import "net/http"

type ContentWriter interface {
	Write(w http.ResponseWriter, httpStatus int, body interface{}) int
	WriteView(w http.ResponseWriter, httpStatus int, view interface{}) int
	WriteError(w http.ResponseWriter, err error) int
}
