package nval

type MetadataBuilder struct {
	metadata map[string]interface{}
}

func (m MetadataBuilder) Add(k string, v interface{}) MetadataBuilder {
	m.metadata[k] = v
	return m
}

func (m MetadataBuilder) Build() map[string]interface{} {
	return m.metadata
}
