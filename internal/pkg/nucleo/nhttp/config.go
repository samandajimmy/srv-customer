package nhttp

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/url"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
)

type ServerConfig struct {
	ListenPort int
	BasePath   string
	BaseUrl    url.URL
	Secure     bool
	TrustProxy bool
	Debug      bool
}

func (s ServerConfig) Validate() error {
	err := s.LoadFromEnv()

	if err != nil {
		return ncore.TraceError("error when load from env", err)
	}

	return validation.ValidateStruct(&s,
		validation.Field(&s.ListenPort, validation.Required, validation.Min(1024), validation.Max(65535)),
		validation.Field(&s.BasePath, validation.Required),
	)
}

func (s *ServerConfig) LoadFromEnv() error {
	// Parse enabled value
	s.ListenPort, _ = nval.ParseInt(os.Getenv("PORT"))
	s.TrustProxy = nval.ParseBooleanFallback(os.Getenv("SERVER_TRUST_PROXY"), false)
	s.Debug = nval.ParseBooleanFallback(os.Getenv("DEBUG"), false)
	s.Secure = nval.ParseBooleanFallback(os.Getenv("SERVER_LISTEN_SECURE"), false)

	// Normalize base path
	basePath := nval.ParseStringFallback(os.Getenv("SERVER_BASE_PATH"), "/")
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}
	s.BasePath = basePath

	// Get base URL
	baseUrlStr, ok := nval.ParseString(os.Getenv("SERVER_HTTP_BASE_URL"))
	if !ok || baseUrlStr == "" {
		log.Warn("server.base_url is not set")
		s.BaseUrl = BuildUrl("localhost", s.ListenPort, s.BasePath)
	} else {
		tmp, err := url.Parse(baseUrlStr)
		if err != nil {
			log.Warn("failed to parse server.base_url = " + baseUrlStr)
			s.BaseUrl = BuildUrl("localhost", s.ListenPort, s.BasePath)
		} else {
			s.BaseUrl = *tmp
			s.BaseUrl.Path = s.BasePath
		}
	}

	return nil
}

func (s *ServerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Unmarshal origin as map interface
	tmp := make(map[string]interface{})
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	// Parse enabled value
	s.ListenPort, _ = nval.ParseInt(tmp["listen_port"])
	s.TrustProxy = nval.ParseBooleanFallback(tmp["trust_proxy"], false)
	s.Debug = nval.ParseBooleanFallback(tmp["debug"], false)
	s.Secure = nval.ParseBooleanFallback(tmp["secure"], false)

	// Normalize base path
	basePath := nval.ParseStringFallback(tmp["base_path"], "/")
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}
	s.BasePath = basePath

	// Get base URL
	baseUrlStr, ok := nval.ParseString(tmp["base_url"])
	if !ok {
		log.Warn("server.base_url is not set")
		s.BaseUrl = BuildUrl("localhost", s.ListenPort, s.BasePath)
	} else {
		tmp, err := url.Parse(baseUrlStr)
		if err != nil {
			log.Warn("failed to parse server.base_url = " + baseUrlStr)
			s.BaseUrl = BuildUrl("localhost", s.ListenPort, s.BasePath)
		} else {
			s.BaseUrl = *tmp
			s.BaseUrl.Path = s.BasePath
		}
	}

	return nil
}

func (s *ServerConfig) GetListenPort() string {
	return fmt.Sprintf(":%d", s.ListenPort)
}

func (s *ServerConfig) GetHttpBaseUrl() string {
	u := s.BaseUrl
	if s.Secure {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	return u.String()
}

func (s *ServerConfig) GetWebSocketBaseUrl() string {
	u := s.BaseUrl
	if s.Secure {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	return u.String()
}

func (s *ServerConfig) GetBasePath() string {
	if s.BasePath == "" {
		return "/"
	}

	if s.BasePath == "/" || !strings.HasSuffix(s.BasePath, "/") {
		return s.BasePath
	}

	return strings.TrimRight(s.BasePath, "/")
}
