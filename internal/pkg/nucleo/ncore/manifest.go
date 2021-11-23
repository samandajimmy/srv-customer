package ncore

import "fmt"

type Manifest struct {
	AppName    string                 `json:"app_name"`
	AppVersion string                 `json:"app_version"`
	Metadata   map[string]interface{} `json:"metadata"`
}

func (m *Manifest) AddMetadata(key string, value interface{}) *Manifest {
	m.Metadata[key] = value
	return m
}

func (m *Manifest) GetMetadata(key string) interface{} {
	v, ok := m.Metadata[key]
	if !ok {
		return nil
	}

	return v
}

func (m *Manifest) GetStringMetadata(key string) string {
	v := m.GetMetadata(key)
	return fmt.Sprintf("%v", v)
}
