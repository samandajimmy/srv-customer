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
	DatabasePort            uint16 `envconfig:"DB_PORT"`
	DatabaseUser            string `envconfig:"DB_USER"`
	DatabasePass            string `envconfig:"DB_PASS"`
	DatabaseName            string `envconfig:"DB_NAME"`
	DatabaseMaxIdleConn     *int   `envconfig:"DB_POOL_MAX_IDLE_CONN"`
	DatabaseMaxOpenConn     *int   `envconfig:"DB_POOL_MAX_OPEN_CONN"`
	DatabaseMaxConnLifetime *int   `envconfig:"DB_POOL_MAX_CONN_LIFETIME"`
	DatabaseBootMigration   bool   `envconfig:"DB_BOOT_MIGRATION"`
}

type ServerConfig struct {
	ListenPort int     `envconfig:"PORT"`
	BasePath   string  `envconfig:"SERVER_BASE_PATH"`
	BaseURL    url.URL `envconfig:"SERVER_HTTP_BASE_URL"`
	Secure     bool    `envconfig:"SERVER_LISTEN_SECURE"`
	TrustProxy string  `envconfig:"SERVER_TRUST_PROXY"`
	Debug      string  `envconfig:"LOG_LEVEL"`
	AssetURL   string
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
	CoreAPIURL         string `envconfig:"CORE_API_URL"`
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
	MinioURL             string `envconfig:"MINIO_URL"`
	MinioSecure          bool   `envconfig:"MINIO_SECURE"`
}

type ClientEndpointConfig struct {
	NotificationServiceURL       string `envconfig:"NOTIFICATION_SERVICE_URL"`
	NotificationServiceAppXid    string `envconfig:"NOTIFICATION_SERVICE_APP_XID"`
	NotificationServiceAppAPIKey string `envconfig:"NOTIFICATION_SERVICE_APP_API_KEY"`
	PdsAPIServiceURL             string `envconfig:"PDS_API_SERVICE_URL"`
}

func (s *ServerConfig) GetListenPort() string {
	return fmt.Sprintf(":%d", s.ListenPort)
}

func (s *ServerConfig) GetHTTPBaseURL() string {
	u := s.BaseURL
	if s.Secure {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	return u.String()
}

func (s *ServerConfig) GetWebSocketBaseURL() string {
	u := s.BaseURL
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
