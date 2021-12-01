package nclient

import (
	"bytes"
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
	var payload = setBodyRequest(body)

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
	fmt.Println(request.Header)
	log.Debugf("Request body: %s", request.Body)
	fmt.Println(request.Body)

	// Do request http with client
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Errorf("Error when request client. err: %s", err)
		return result, ncore.TraceError(err)
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error while reading the response bytes:", err)
		}
		log.Debugf(string([]byte(body)))
		fmt.Println(string([]byte(body)))
	}

	// Set result
	result = resp

	return result, nil
}

func setBodyRequest(data map[string]string) *bytes.Buffer {
	// Set param for body request
	var param = url.Values{}
	for key, value := range data {
		param.Set(key, value)
	}

	return bytes.NewBufferString(param.Encode())
}

func setHeaderRequest(request *http.Request, data map[string]string) *http.Request {
	// setHeaderRequest
	for key, value := range data {
		request.Header.Add(key, value)
	}
	return request
}
