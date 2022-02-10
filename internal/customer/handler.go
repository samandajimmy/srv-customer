package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nredis"
	"time"
)

func NewHandler(core *ncore.Core, config *Config) (*Handler, error) {
	// Init repository
	repoInternal, err := NewRepository(&config.DatabaseConfig)
	if err != nil {
		return nil, ncore.TraceError("failed to internal service", err)
	}

	// Init additional repository
	repoExternal, err := NewRepositoryExternal(&config.DatabaseExternal)
	if err != nil {
		return nil, ncore.TraceError("failed to external service", err)
	}

	// Init Redis
	redisConfig := config.RedisConfig
	redis := nredis.NewNucleoRedis(
		redisConfig.RedisScheme,
		redisConfig.RedisHost,
		redisConfig.RedisPort,
		redisConfig.RedisPass,
	)

	h := Handler{
		startedAt:    time.Now(),
		Core:         core,
		config:       config,
		repo:         repoInternal,
		repoExternal: repoExternal,
		redis:        redis,
	}

	return &h, nil
}

type Handler struct {
	// Service
	*ncore.Core
	config       *Config
	repo         *Repository         // pgsql
	repoExternal *RepositoryExternal // mysql
	// Redis
	redis *nredis.Redis
	// Metadata
	startedAt time.Time
}

type HandlerMap struct {
	// TODO Middleware
	//Middlewares  *handler.Middlewares
	Auth         *Auth
	Common       *Common
	Customer     *Customer
	Verification *Verification
}

func RegisterHandler(manifest ncore.Manifest, h *Handler) *HandlerMap {
	return &HandlerMap{
		Common:       NewCommon(time.Now(), manifest),
		Auth:         NewAuth(h),
		Customer:     NewCustomer(h),
		Verification: NewVerification(h),
	}
}
