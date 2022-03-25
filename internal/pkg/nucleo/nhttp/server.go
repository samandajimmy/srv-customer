package nhttp

import "fmt"

// ListenPort returns value that can satisfy addr value when calling http.ListenAndServe
func ListenPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
