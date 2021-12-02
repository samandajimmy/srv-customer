package nclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
)

type Nclient struct {
	ChannelId string
	ClientId  string
	Client    http.Client
	BaseUrl   string
}

type ResponseSwitching struct {
	Message string `json:"data"`
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

func (c *Nclient) PostData(endpoint string, body map[string]string, header map[string]string) (*http.Response, error) {
	var result *http.Response
	var payload *bytes.Buffer

	// Get body request
	payload = getBodyRequest(header, body)

	// Make request http
	endPoint := c.BaseUrl + endpoint
	request, err := http.NewRequest("POST", endPoint, payload)
	if err != nil {
		log.Errorf("Error when make new request. err: %s", err)
		return result, ncore.TraceError(err)
	}

	// Set header
	request = setHeaderRequest(request, header)
	log.Debugf("Request header: %s", request.Header)
	log.Debugf("Request body: %s", request.Body)

	// Do request http with client
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Errorf("Error when request client. err: %s", err)
		return result, ncore.TraceError(err)
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

func GetResponseData(response *http.Response) (string, error) {
	defer response.Body.Close()
	var Response *ResponseSwitching
	err := json.NewDecoder(response.Body).Decode(&Response)
	if err != nil {
		log.Errorf("Error while reading the response bytes:", err)
		return Response.Message, err
	}
	return Response.Message, nil
}

func getBodyRequest(header map[string]string, body map[string]string) *bytes.Buffer {
	var payload *bytes.Buffer
	switch header["Content-Type"] {
	case "application/json":
		payload = setBodyApplicationJSON(body)
		break
	case "application/x-www-form-urlencoded":
		payload = setBodyUrlEncoded(body)
		break
	default:
		payload = setBodyApplicationJSON(body)
		break
	}
	return payload
}

func setBodyUrlEncoded(data map[string]string) *bytes.Buffer {
	var param = url.Values{}
	for key, value := range data {
		param.Set(key, value)
	}

	return bytes.NewBufferString(param.Encode())
}

func setBodyApplicationJSON(data map[string]string) *bytes.Buffer {
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
