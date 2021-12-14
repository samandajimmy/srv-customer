package contract

import (
	"encoding/hex"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type Config struct {
	Server      nhttp.ServerConfig
	Client      ClientConfig
	DataSources DataSourcesConfig
	CORS        nhttp.CORSConfig
	SMTP        SMTPConfig
	CorePDS     CorePDSConfig
	Redis       RedisConfig
}

func (c *Config) LoadFromEnv() {
	// Set config server
	port, _ := nval.ParseInt(os.Getenv("PORT"))
	c.Server.ListenPort = port
	if c.Server.ListenPort == 0 {
		c.Server.ListenPort = 3000
	}

	if c.Server.BasePath == "" {
		c.Server.BasePath = os.Getenv("SERVER_BASE_PATH")
	}

	c.Server.Secure = nval.ParseBooleanFallback(os.Getenv("SERVER_LISTEN_SECURE"), true)
	c.Server.TrustProxy = nval.ParseBooleanFallback(os.Getenv("SERVER_TRUST_PROXY"), true)
	c.Server.Debug = nval.ParseBooleanFallback(os.Getenv("DEBUG"), false)

	// Set config client
	c.Client.ClientID = nval.ParseStringFallback(os.Getenv("CLIENT_ID"), "")
	c.Client.ClientSecret = nval.ParseStringFallback(os.Getenv("CLIENT_SECRET"), "")
	c.Client.JWTExpired = nval.ParseInt64Fallback(os.Getenv("JWT_EXP"), 3600)
	c.Client.JWTKey = nval.ParseStringFallback(os.Getenv("JWT_KEY"), nval.RandStringBytes(78))

	// Set config data resource internal
	c.DataSources.DBInternal = nsql.Config{
		Driver:          os.Getenv("INTERNAL_DB_DRIVER"),
		Host:            os.Getenv("INTERNAL_DB_HOST"),
		Port:            os.Getenv("INTERNAL_DB_PORT"),
		Username:        os.Getenv("INTERNAL_DB_USER"),
		Password:        os.Getenv("INTERNAL_DB_PASS"),
		Database:        os.Getenv("INTERNAL_DB_NAME"),
		MaxIdleConn:     nsql.NewInt(10),
		MaxOpenConn:     nsql.NewInt(10),
		MaxConnLifetime: nsql.NewInt(1),
	}
	// DB EXTERNAL
	c.DataSources.DBExternal = nsql.Config{
		Driver:          os.Getenv("EXTERNAL_DB_DRIVER"),
		Host:            os.Getenv("EXTERNAL_DB_HOST"),
		Port:            os.Getenv("EXTERNAL_DB_PORT"),
		Username:        os.Getenv("EXTERNAL_DB_USER"),
		Password:        os.Getenv("EXTERNAL_DB_PASS"),
		Database:        os.Getenv("EXTERNAL_DB_NAME"),
		MaxIdleConn:     nsql.NewInt(10),
		MaxOpenConn:     nsql.NewInt(10),
		MaxConnLifetime: nsql.NewInt(1),
	}

	// If password is hex
	password, err := hex.DecodeString(fmt.Sprintf("%v", os.Getenv("EXTERNAL_DB_PASS")))
	c.DataSources.DBExternal.Password = string(password)
	if err != nil {
		log.Errorf("Error parsing hex env for password: %v", err)
		c.DataSources.DBExternal.Password = os.Getenv("EXTERNAL_DB_PASS")
	}

	// Load cors
	corsEnabled := nval.ParseBooleanFallback(os.Getenv("CORS_ENABLED"), false)
	if corsEnabled {
		c.CORS = nhttp.CORSConfig{
			Enabled:        true,
			Origins:        nval.ParseStringArrayFallback(os.Getenv("CORS_ORIGINS"), []string{"*"}),
			AllowedHeaders: nval.ParseStringArrayFallback(os.Getenv("CORS_ALLOWED_HEADERS"), []string{"*"}),
			AllowedMethods: nval.ParseStringArrayFallback(os.Getenv("CORS_ALLOWED_METHODS"), []string{http.MethodGet,
				http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions}),
		}
	}

	// Load smtp config
	smtpPort, _ := nval.ParseInt(os.Getenv("SMTP_PORT"))
	c.SMTP = SMTPConfig{
		Host:     nval.ParseStringFallback(os.Getenv("SMTP_HOST"), ""),
		Port:     nval.ParseIntFallback(smtpPort, 587),
		Username: nval.ParseStringFallback(os.Getenv("SMTP_USERNAME"), ""),
		Password: nval.ParseStringFallback(os.Getenv("SMTP_PASSWORD"), ""),
	}

	// Load PDS CORE API Config
	c.CorePDS = CorePDSConfig{
		CoreApiUrl:         nval.ParseStringFallback(os.Getenv("CORE_API_URL"), ""),
		CoreOauthUsername:  nval.ParseStringFallback(os.Getenv("CORE_OAUTH_USERNAME"), ""),
		CoreOauthPassword:  nval.ParseStringFallback(os.Getenv("CORE_OAUTH_PASSWORD"), ""),
		CoreOauthGrantType: nval.ParseStringFallback(os.Getenv("CORE_OAUTH_GRANT_TYPE"), ""),
		CoreAuthorization:  nval.ParseStringFallback(os.Getenv("CORE_AUTHORIZATION"), ""),
		CoreClientId:       nval.ParseStringFallback(os.Getenv("CORE_CLIENT_ID"), ""),
	}

	// Load Redis Config
	c.Redis = RedisConfig{
		RedisScheme: nval.ParseStringFallback(os.Getenv("REDIS_SCHEME"), "tcp"),
		RedisHost:   nval.ParseStringFallback(os.Getenv("REDIS_HOST"), "localhost"),
		RedisPort:   nval.ParseStringFallback(os.Getenv("REDIS_PORT"), "6379"),
		RedisPass:   nval.ParseStringFallback(os.Getenv("REDIS_PASS"), ""),
		RedisExpiry: nval.ParseInt64Fallback(os.Getenv("REDIS_EXPIRY"), 0),
	}
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Server),
		validation.Field(&c.DataSources),
	)
}

type DataSourcesConfig struct {
	DBInternal nsql.Config
	DBExternal nsql.Config
}

func (c DataSourcesConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DBInternal, validation.Required),
	)
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	JWTExpired   int64
	JWTKey       string
}

func (c ClientConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClientID, validation.Required),
		validation.Field(&c.ClientSecret, validation.Required),
	)
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (c SMTPConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
	)
}

type CorePDSConfig struct {
	CoreApiUrl         string
	CoreOauthUsername  string
	CoreOauthPassword  string
	CoreOauthGrantType string
	CoreAuthorization  string
	CoreClientId       string
}

func (c CorePDSConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CoreApiUrl, validation.Required),
		validation.Field(&c.CoreOauthUsername, validation.Required),
		validation.Field(&c.CoreOauthPassword, validation.Required),
		validation.Field(&c.CoreOauthGrantType, validation.Required),
		validation.Field(&c.CoreAuthorization, validation.Required),
		validation.Field(&c.CoreClientId, validation.Required),
	)
}

type RedisConfig struct {
	RedisScheme string
	RedisHost   string
	RedisPort   string
	RedisPass   string
	RedisExpiry int64
}

func (c RedisConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.RedisScheme, validation.Required),
		validation.Field(&c.RedisHost, validation.Required),
		validation.Field(&c.RedisPort, validation.Required),
	)
}
