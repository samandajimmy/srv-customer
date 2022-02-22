package nhttp

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(args ...RouterOptions) *Router {
	// Get options
	var options RouterOptions
	if len(args) > 0 {
		options = args[0]
	} else {
		options = RouterOptions{
			LogRequest: true,
			Debug:      false,
			TrustProxy: false,
		}
	}

	// If content writer is not set, then set to json
	if options.ContentWriter == nil {
		options.ContentWriter = JSONContentWriter{Debug: options.Debug}
	}

	// If router is not set then initiate a new router
	var muxRouter *mux.Router
	if options.Router == nil {
		muxRouter = mux.NewRouter()
	} else {
		muxRouter = options.Router
	}

	// Init router
	router := Router{
		Router:        muxRouter,
		Debug:         options.Debug,
		LogRequest:    options.LogRequest,
		contentWriter: options.ContentWriter,
	}

	// If Log Request is enabled, then set logging
	if options.LogRequest {
		// Init pre log request middleware
		router.handleCaptureRequestMetadata = NewCaptureRequestMetadataHandler(options.TrustProxy)
		// Set Middleware
		router.Use(router.handleCaptureRequestMetadata)
	}

	return &router
}

type RouterOptions struct {
	Router        *mux.Router
	LogRequest    bool
	Debug         bool
	TrustProxy    bool
	ContentWriter ContentWriter
}

type Router struct {
	*mux.Router
	Debug      bool
	LogRequest bool
	// Private
	contentWriter                ContentWriter
	handleCaptureRequestMetadata mux.MiddlewareFunc
}

func (r *Router) HandleFunc(handlerFn HandlerFunc) http.Handler {
	return NewHandler(handlerFn).SetWriter(r.contentWriter)
}

func (r *Router) Handle(method, path string, handler http.Handler, nextHandlers ...http.Handler) *mux.Route {
	var h http.Handler

	// Determine handler
	if len(nextHandlers) == 0 {
		h = handler
	} else {
		// Init chain builder
		hArr := append([]http.Handler{handler}, nextHandlers...)
		chainBuilder := NewHandlerChain(hArr...)

		// Create a chain middleware function
		h = chainBuilder.Build()
	}

	// If log request is enabled, then wrap with HandleLogRequest
	if r.LogRequest {
		h = HandleLogRequest(h)
	}

	return r.NewRoute().Path(path).Handler(h).Methods(method)
}

func (r *Router) RESTSubRouter(path string) *Router {
	// Init sub router
	s := NewRouter(RouterOptions{
		Router:        r.PathPrefix(path).Subrouter(),
		LogRequest:    r.LogRequest,
		Debug:         r.Debug,
		ContentWriter: JSONContentWriter{Debug: r.Debug},
	})

	s.HandleErrorNotFound(NewHandler(HandleErrorNotFound).SetWriter(s.contentWriter))
	s.HandleErrorMethodNotAllowed(NewHandler(HandleErrorMethodNotAllowed).SetWriter(s.contentWriter))

	return s
}

func (r *Router) HandleErrorNotFound(h http.Handler) {
	// Check HandleLogRequest options
	if r.handleCaptureRequestMetadata != nil {
		r.NotFoundHandler = HandleLogRequest(h)
	} else {
		r.NotFoundHandler = h
	}
}

func (r *Router) HandleErrorMethodNotAllowed(h http.Handler) {
	// Check HandleLogRequest options
	if r.handleCaptureRequestMetadata != nil {
		r.MethodNotAllowedHandler = HandleLogRequest(h)
	} else {
		r.MethodNotAllowedHandler = h
	}
}
