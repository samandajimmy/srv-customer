package customer

import (
	"fmt"
	"net/url"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"strings"
)

type Config struct {
	ServerConfig
	DatabaseConfig
	DatabaseExternal DatabaseConfig
	SMTPConfig
	EmailConfig
	MinioConfig
	ClientConfig
	RedisConfig
	CoreSwitchingConfig
	ClientEndpointConfig
	CORS nhttp.CORSConfig
}

type DatabaseConfig struct {
	DatabaseDriver          string `envconfig:"DB_DRIVER"`
	DatabaseHost            string `envconfig:"DB_HOST"`
	DatabasePort            string `envconfig:"DB_PORT"`
	DatabaseUser            string `envconfig:"DB_USER"`
	DatabasePass            string `envconfig:"DB_PASS"`
	DatabaseName            string `envconfig:"DB_NAME"`
	DatabaseMaxIdleConn     *int   `envconfig:"DB_POOL_MAX_IDLE_CONN"`
	DatabaseMaxOpenConn     *int   `envconfig:"DB_POOL_MAX_OPEN_CONN"`
	DatabaseMaxConnLifetime *int   `envconfig:"DB_POOL_MAX_CONN_LIFETIME"`
}

type ServerConfig struct {
	ListenPort int     `envconfig:"PORT"`
	BasePath   string  `envconfig:"SERVER_BASE_PATH"`
	BaseUrl    url.URL `envconfig:"SERVER_HTTP_BASE_URL"`
	Secure     bool    `envconfig:"SERVER_LISTEN_SECURE"`
	TrustProxy string  `envconfig:"SERVER_TRUST_PROXY"`
	Debug      string  `envconfig:"LOG_LEVEL"`
	AssetUrl   string
}

type SMTPConfig struct {
	SMTPHost     string `envconfig:"SMTP_HOST"`
	SMTPPort     string `envconfig:"SMTP_PORT"`
	SMTPUsername string `envconfig:"SMTP_USERNAME"`
	SMTPPassword string `envconfig:"SMTP_PASSWORD"`
}

type EmailConfig struct {
	PdsEmailFrom     string `envconfig:"PDS_EMAIL_FROM"`
	PdsEmailFromName string `envconfig:"PDS_EMAIL_FROM_NAME"`
}

type ClientConfig struct {
	ClientID     string `envconfig:"CLIENT_ID"`
	ClientSecret string `envconfig:"CLIENT_SECRET"`
	JWTExpiry    int64  `envconfig:"JWT_EXP"`
	JWTKey       string `envconfig:"JWT_KEY"`
}

type CoreSwitchingConfig struct {
	CoreApiURL         string `envconfig:"CORE_API_URL"`
	CoreOauthUsername  string `envconfig:"CORE_OAUTH_USERNAME"`
	CoreOauthPassword  string `envconfig:"CORE_OAUTH_PASSWORD"`
	CoreOauthGrantType string `envconfig:"CORE_OAUTH_GRANT_TYPE"`
	CoreAuthorization  string `envconfig:"CORE_AUTHORIZATION"`
	CoreClientID       string `envconfig:"CORE_CLIENT_ID"`
}

type RedisConfig struct {
	RedisScheme string `envconfig:"REDIS_SCHEME"`
	RedisHost   string `envconfig:"REDIS_HOST"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	RedisPass   string `envconfig:"REDIS_PASS"`
	RedisExpiry int64  `envconfig:"REDIS_EXPIRY"`
}

type MinioConfig struct {
	MinioAccessKeyID     string `envconfig:"MINIO_ACCESS_KEY_ID"`
	MinioSecretAccessKey string `envconfig:"MINIO_SECRET_ACCESS_KEY"`
	MinioBucket          string `envconfig:"MINIO_BUCKET"`
	MinioEndpoint        string `envconfig:"MINIO_ENDPOINT"`
	MinioUrl             string `envconfig:"MINIO_URL"`
	MinioSecure          bool   `envconfig:"MINIO_SECURE"`
}

type ClientEndpointConfig struct {
	NotificationServiceUrl string `envconfig:"NOTIFICATION_SERVICE_URL"`
	PdsApiServiceUrl       string `envconfig:"PDS_API_SERVICE_URL"`
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
