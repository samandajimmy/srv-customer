package nhttp

type ResponseFlag = int

const (
	EndRequest = ResponseFlag(iota)
	ContinueRequest
)

const (
	SuccessCode    = "OK"
	SuccessMessage = "Success"
)

func OK() *Response {
	return &Response{
		Success:      true,
		Code:         SuccessCode,
		Message:      SuccessMessage,
		Data:         nil,
		Header:       make(map[string]string),
		responseFlag: EndRequest,
	}
}

func Success() *Response {
	return &Response{
		Success:      true,
		Code:         SuccessCode,
		Message:      SuccessMessage,
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
