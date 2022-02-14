package contract

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Config struct {
	Server         nhttp.ServerConfig
	DataSources    DataSourcesConfig
	Client         ClientConfig
	CORS           nhttp.CORSConfig
	SMTP           SMTPConfig
	CorePDS        CorePDSConfig
	Redis          RedisConfig
	ClientEndpoint ClientEndpointConfig
	Email          EmailConfig
}

type ClientEndpointConfig struct {
	NotificationServiceUrl string
	PdsApiServiceUrl       string
}

type DataSourcesConfig struct {
	DBInternal nsql.Config
	DBExternal nsql.Config
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	JWTExpired   int64
	JWTKey       string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type CorePDSConfig struct {
	CoreApiUrl         string
	CoreOauthUsername  string
	CoreOauthPassword  string
	CoreOauthGrantType string
	CoreAuthorization  string
	CoreClientId       string
}

type RedisConfig struct {
	RedisScheme string
	RedisHost   string
	RedisPort   string
	RedisPass   string
	RedisExpiry int64
}

type EmailConfig struct {
	PdsEmailFrom     string
	PdsEmailFromName string
}
