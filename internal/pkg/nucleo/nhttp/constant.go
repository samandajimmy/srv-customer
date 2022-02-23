package nhttp

const (
	// Header keys

	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"

	// Map keys

	MetadataKey       = "metadata"
	HTTPStatusRespKey = "http_status"

	NotApplicable = "N/A"
)

// Context Key

type ContextKey uint8

const (
	_ ContextKey = iota
	RequestIDContextKey
	HTTPStatusRespContextKey
	RequestMetadataContextKey
)
