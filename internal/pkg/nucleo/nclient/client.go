package nclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type Nclient struct {
	ChannelId string
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

func NewNucleoClient(channelId string, clientId string, baseUrl string) *Nclient {
	return &Nclient{
		ChannelId: channelId,
		ClientId:  clientId,
		Client:    http.Client{},
		BaseUrl:   baseUrl,
	}
}

func (c *Nclient) PostData(endpoint string, body map[string]interface{}, header map[string]string) (*http.Response, error) {
	var result *http.Response
	var payload *bytes.Buffer

	// Get body request
	payload, header = getBodyRequest(header, body)

	// Make request http
	baseUrlWithEndpoint := c.BaseUrl + endpoint
	log.Debugf("PostData.Endpoint %s", baseUrlWithEndpoint)
	request, err := http.NewRequest("POST", baseUrlWithEndpoint, payload)
	if err != nil {
		log.Errorf("Error when make new request. err: %s", err)
		return result, ncore.TraceError("Error when make new request", err)
	}

	// Set header
	request = setHeaderRequest(request, header)
	log.Debugf("Request header: %s", request.Header)
	log.Debugf("Request body: %s", request.Body)

	// Do request http with client
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Errorf("Error when request client. err: %s", err)
		return result, ncore.TraceError("Error when request client", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(GetResponseString(resp))
	}

	// Set result
	result = resp

	return result, nil
}

func GetResponseString(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error while reading the response bytes:", err)
		return ""
	}
	log.Debugf(string(body))
	return fmt.Sprintf(string(body))
}

func getBodyRequest(header map[string]string, body map[string]interface{}) (*bytes.Buffer, map[string]string) {
	var payload *bytes.Buffer
	switch header["Content-Type"] {
	case "application/json":
		payload = setBodyApplicationJSON(body)
		break
	case "application/x-www-form-urlencoded":
		payload = setBodyUrlEncoded(body)
		break
	case "multipart/form-data":
		result, contentType := setBodyFormData(body)
		header["Content-Type"] = contentType
		payload = result
		break
	default:
		payload = setBodyApplicationJSON(body)
		break
	}
	return payload, header
}

func setBodyUrlEncoded(data map[string]interface{}) *bytes.Buffer {
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
