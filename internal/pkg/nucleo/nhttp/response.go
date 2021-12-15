package nhttp

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type ResponseFlag = int

const (
	EndRequest = ResponseFlag(iota)
	ContinueRequest
	ViewRequest
)

func NewResponse() *Response {
	return &Response{
		Header:       make(map[string]string),
		responseFlag: EndRequest,
	}
}

func OK() *Response {
	return &Response{
		Success:      true,
		Code:         ncore.Success.Code,
		Message:      ncore.Success.Message,
		Data:         nil,
		Header:       make(map[string]string),
		responseFlag: EndRequest,
	}
}

func Success() *Response {
	return &Response{
		Success:      true,
		Code:         ncore.Success.Code,
		Message:      ncore.Success.Message,
		Header:       make(map[string]string),
		responseFlag: EndRequest,
	}
}

func View() *Response {
	return &Response{
		Success:      true,
		Code:         ncore.Success.Code,
		Message:      ncore.Success.Message,
		Header:       make(map[string]string),
		Data:         "",
		responseFlag: ViewRequest,
	}
}

func BadRequest(data interface{}) *Response {
	return &Response{
		Success:      true,
		Code:         BadRequestError.Code,
		Message:      BadRequestError.Message,
		Data:         data,
		Header:       make(map[string]string),
		responseFlag: EndRequest,
	}
}

func UnprocessableEntity(data interface{}) *Response {
	return &Response{
		Success:      false,
		Code:         UnprocessableEntityError.Code,
		Message:      UnprocessableEntityError.Message,
		Data:         data,
		Header:       make(map[string]string),
		responseFlag: EndRequest,
	}
}

func Continue() *Response {
	return &Response{
		responseFlag: ContinueRequest,
	}
}

// Response represents base response structure if response is being handled properly
type Response struct {
	Success bool              `json:"success"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data"`
	Header  map[string]string `json:"-"`
	// Set response flow
	responseFlag ResponseFlag
}

func (r *Response) AddHeader(key string, value string) *Response {
	r.Header[key] = value
	return r
}

func (r *Response) SetCode(code string) *Response {
	r.Code = code
	return r
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) SetSuccess(success bool) *Response {
	r.Success = success
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}
