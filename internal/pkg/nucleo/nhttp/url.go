package nhttp

import (
	"fmt"
	"net/url"
)

func BuildURL(host string, port int, basePath string) url.URL {
	if port != 443 && port != 80 {
		host = fmt.Sprintf("%s:%d", host, port)
	}

	return url.URL{
		Host: host,
		Path: basePath,
	}
}
