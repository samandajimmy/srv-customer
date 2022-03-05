package nhttp

const (
	pkgNamespace = "nhttp"
)

const (
	// Header keys

	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"

	// Map keys

	MetadataKey = "metadata"

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
