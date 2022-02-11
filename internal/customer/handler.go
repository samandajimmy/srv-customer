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
		StartedAt:    time.Now(),
		Core:         core,
		Config:       config,
		Redis:        redis,
		Repo:         repoInternal,
		RepoExternal: repoExternal,
	}

	return &h, nil
}

type Handler struct {
	// Repo
	Repo         *Repository         // pgsql
	RepoExternal *RepositoryExternal // mysql
	// Redis
	Redis *nredis.Redis
	// Metadata
	*ncore.Core
	StartedAt time.Time
	Config    *Config
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
