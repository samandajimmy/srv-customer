package customer

import (
	"github.com/nbs-go/errx"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nredis"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ns3"
	"time"
)

//nolint:funlen
func NewHandler(core *ncore.Core, config *Config) (*Handler, error) {
	// Init repository
	repoInternal, err := NewRepository(&config.DatabaseConfig)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Init additional repository
	repoExternal, err := NewRepositoryExternal(&config.DatabaseExternal)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Init Redis
	redisConfig := config.RedisConfig
	redis := nredis.NewNucleoRedis(
		redisConfig.RedisScheme,
		redisConfig.RedisHost,
		redisConfig.RedisPort,
		redisConfig.RedisPass,
	)

	httpClient := http.Client{
		Timeout: time.Minute,
	}

	// Init base client
	client := &nclient.Nclient{
		Client: httpClient,
	}

	// Init PDS API client config
	pdsAPIClient := &nclient.Nclient{
		Client:  httpClient,
		BaseUrl: config.PdsAPIServiceURL,
	}

	// Initialize minio client object.
	minioClient, err := ns3.NewMinio(ns3.MinioOpt{
		Endpoint:        config.MinioEndpoint,
		AccessKeyID:     config.MinioAccessKeyID,
		SecretAccessKey: config.MinioSecretAccessKey,
		UseSSL:          config.MinioSecure,
		BucketName:      config.MinioBucket,
	})
	if err != nil {
		return nil, errx.Trace(err)
	}

	h := Handler{
		StartedAt:    time.Now(),
		Core:         core,
		Config:       config,
		Redis:        redis,
		Client:       client,
		PdsAPIClient: pdsAPIClient,
		Minio:        minioClient,
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
	// Client
	Client       *nclient.Nclient
	PdsAPIClient *nclient.Nclient
	Minio        *ns3.Minio
	// Metadata
	*ncore.Core
	StartedAt time.Time
	Config    *Config
}

type Controllers struct {
	Middlewares  *Middlewares
	Auth         *Auth
	Common       *CommonController
	Customer     *Customer
	Asset        *Asset
	Verification *Verification
	Profile      *ProfileController
}

func NewControllers(manifest ncore.Manifest, h *Handler) *Controllers {
	return &Controllers{
		Common:       NewCommonController(time.Now(), manifest),
		Asset:        NewAsset(h),
		Middlewares:  NewMiddlewares(h),
		Auth:         NewAuth(h),
		Customer:     NewCustomer(h),
		Verification: NewVerification(h),
		Profile:      NewProfileController(h),
	}
}
