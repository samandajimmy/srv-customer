package nclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbs-go/errx"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type Nclient struct {
	ChannelID string
	ClientId  string
	Client    http.Client
	BaseUrl   string
}

type ResponseSwitching struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	Message      string `json:"data"`
}

type ResponsePdsAPI struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

var log = nlogger.Get()

func (c *Nclient) PostData(endpoint string, body map[string]interface{}, header map[string]string) (*http.Response, error) {
	var payload *bytes.Buffer

	// Get body request
	payload, header = getBodyRequest(header, body)

	// Make request http
	baseURLWithEndpoint := c.BaseUrl + endpoint
	log.Debugf("PostData.Endpoint %s", baseURLWithEndpoint)
	request, err := http.NewRequest("POST", baseURLWithEndpoint, payload)
	if err != nil {
		log.Errorf("Error when make new request. err: %s", err)
		return nil, errx.Trace(err)
	}

	// Set header
	request = setHeaderRequest(request, header)
	log.Debugf("Request header: %s", request.Header)
	log.Debugf("Request body: %s", request.Body)

	// Do request http with client
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Errorf("Error when request client. err: %s", err)
		return nil, errx.Trace(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(GetResponseString(resp))
	}

	return resp, nil
}

func GetResponseString(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error while reading the response bytes:", err)
		return ""
	}
	log.Debugf(string(body))
	return string(body)
}

func getBodyRequest(header map[string]string, body map[string]interface{}) (*bytes.Buffer, map[string]string) {
	var payload *bytes.Buffer
	switch header["Content-Type"] {
	case "application/json":
		payload = setBodyApplicationJSON(body)
	case "application/x-www-form-urlencoded":
		payload = setBodyURLEncoded(body)
	case "multipart/form-data":
		result, contentType := setBodyFormData(body)
		header["Content-Type"] = contentType
		payload = result
	default:
		payload = setBodyApplicationJSON(body)
	}
	return payload, header
}

func setBodyURLEncoded(data map[string]interface{}) *bytes.Buffer {
	var param = url.Values{}
	for key, value := range data {
		param.Set(key, nval.ParseStringFallback(value, ""))
	}

	return bytes.NewBufferString(param.Encode())
}

func setBodyFormData(data map[string]interface{}) (*bytes.Buffer, string) {
	result := &bytes.Buffer{}
	writer := multipart.NewWriter(result)
	for key, value := range data {
		_ = writer.WriteField(key, nval.ParseStringFallback(value, ""))
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	contentType := writer.FormDataContentType()

	return result, contentType
}

func setBodyApplicationJSON(data map[string]interface{}) *bytes.Buffer {
	// Set param for body request
	jsonValue, _ := json.Marshal(data)

	return bytes.NewBuffer(jsonValue)
}

func setHeaderRequest(request *http.Request, data map[string]string) *http.Request {
	// setHeaderRequest
	for key, value := range data {
		request.Header.Add(key, value)
	}
	return request
}
