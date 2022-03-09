package nhttp

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"time"
)

type RequestMetadata struct {
	ClientIP  string
	StartedAt time.Time
}

func NewCaptureRequestMetadataHandler(trustProxy bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture client IP Client IP
			clientIP := GetClientIP(r, trustProxy)
			startedAt := time.Now()

			// Set to context value
			ctx := context.WithValue(r.Context(), RequestMetadataContextKey, RequestMetadata{
				ClientIP:  clientIP,
				StartedAt: startedAt,
			})

			// Set request id value
			reqID, _ := uuid.NewUUID()
			ctx = context.WithValue(ctx, RequestIDContextKey, reqID.String())

			// Continue
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HandleLogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve http
		h.ServeHTTP(w, r)

		// Get context
		ctx := r.Context()

		// Get Log Request Metadata
		var elapsedTime, clientIP string
		metadata, ok := ctx.Value(RequestMetadataContextKey).(RequestMetadata)
		if !ok {
			elapsedTime = "N/A"
			clientIP = "N/A"
		} else {
			elapsedTime = time.Since(metadata.StartedAt).String()
			clientIP = metadata.ClientIP
		}

		// Get httpStatus
		httpStatus, ok := ctx.Value(HTTPStatusRespContextKey).(int)
		if !ok {
			httpStatus = -1
		}

		log.Infof("Endpoint: %s %s, RespHTTPStatus: %d, ElapsedTime: %s, ClientIP: %s", r.Method,
			html.EscapeString(r.URL.Path), httpStatus, elapsedTime, clientIP)
	})
}
