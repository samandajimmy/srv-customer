package ncore

import "fmt"

type Manifest struct {
	AppName        string                 `json:"appName"`
	AppVersion     string                 `json:"appVersion"`
	BuildSignature string                 `json:"buildSignature"`
	Metadata       map[string]interface{} `json:"metadata"`
}

func NewManifest(appName, appVersion, buildSignature string) Manifest {
	return Manifest{
		AppName:        appName,
		AppVersion:     appVersion,
		BuildSignature: buildSignature,
		Metadata:       make(map[string]interface{}),
	}
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
