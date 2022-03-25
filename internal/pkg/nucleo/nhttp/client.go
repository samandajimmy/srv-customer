package nhttp

import (
	logOption "github.com/nbs-go/nlogger/v2/option"
	"net"
	"net/http"
	"strings"

	"github.com/nbs-go/nlogger/v2"
)

var log nlogger.Logger

/// init must be done at the first file in package that is sorted alphabetically
func init() {
	log = nlogger.Get()
}

func GetClientIP(req *http.Request, trustProxy bool) string {
	// Retrieve client IP from x-real-ip header
	tmp := req.Header.Get("x-real-ip")

	// If empty, retrieve client IP from x-forwarded-for header
	if tmp == "" && trustProxy {
		tmp = req.Header.Get("x-forwarded-for")
	}

	// If still empty, then retrieve from req.RemoteAddr
	if tmp == "" {
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			log.Error("unable to get Client IP address", logOption.Error(err))
		}
		tmp = host
	}

	// Split with comma
	remoteAddresses := strings.Split(tmp, ",")
	return strings.TrimSpace(remoteAddresses[0])
}
