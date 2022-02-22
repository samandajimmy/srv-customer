package ncore

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func NewResponseMap() ResponseMap {
	return ResponseMap{
		responses: make(map[string]Response),
	}
}

type ResponseMap struct {
	responses map[string]Response
}

func (m *ResponseMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Unmarshal origin as map of temporary response
	tmp := make(map[string]Response)
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	// Set response codes
	for k, v := range tmp {
		v.Code = k
		m.responses[k] = v
	}

	return nil
}

func (m *ResponseMap) GetSuccess(code string) Response {
	resp, ok := m.responses[code]

	// If not found, then return standard success response
	if !ok {
		return Success
	}

	// If found, but Success if false, then return standard success response
	if !resp.Success {
		return Success
	}

	return resp
}

func (m *ResponseMap) GetError(code string) *Response {
	resp, ok := m.responses[code]

	// If not found, then return standard success response
	if !ok {
		return InternalError
	}

	// If found, but Success if true, then return standard success response
	if resp.Success {
		return InternalError
	}

	return &resp
}

func (m *ResponseMap) Add(resp Response) {
	m.responses[resp.Code] = resp
}

func loadResponseMap(responseMapFile string) (*ResponseMap, error) {
	// Init response map
	respMap := NewResponseMap()

	// If response map file path is set, load from file
	if responseMapFile != "" {
		// Load file
		file, err := ioutil.ReadFile(responseMapFile)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to read ResponseMap file (responseMapFile = %s)"+wrappedErrorFmt, namespace, responseMapFile, err)
		}

		// Parse YAML file
		err = yaml.Unmarshal(file, &respMap)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse ResponseMap file. (responseMapFile = %s)"+wrappedErrorFmt, namespace, responseMapFile, err)
		}
	}

	// Set standard responses
	respMap.Add(Success)
	respMap.Add(*InternalError)

	return &respMap, nil
}
