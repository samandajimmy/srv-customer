package customer

import (
	"context"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nredis"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ns3"
)

type Service struct {
	config       *Config
	ctx          context.Context
	repo         *RepositoryContext
	repoExternal *RepositoryContext
	log          nlogger.Logger
	redis        *nredis.Redis
	minio        *ns3.Minio
	client       *nclient.Nclient
	pdsClient    *nclient.Nclient
}

func (h Handler) NewService(ctx context.Context) *Service {
	svc := Service{
		config:       h.Config,
		client:       h.Client,
		pdsClient:    h.PdsAPIClient,
		redis:        h.Redis,
		minio:        h.Minio,
		repo:         h.Repo.WithContext(ctx),
		repoExternal: h.RepoExternal.WithContext(ctx),
		ctx:          ctx,
		log:          nlogger.NewChild(nlogger.WithNamespace("service"), nlogger.Context(ctx)),
	}

	return &svc
}

func (s *Service) Close() {
	// Close database connection to free pool
	err := s.repo.conn.Close()
	if err != nil {
		s.log.Error("Failed to close connection", nlogger.Error(err))
	}

	// Close database external connection to free pool
	err = s.repoExternal.conn.Close()
	if err != nil {
		s.log.Error("Failed to close external connection", nlogger.Error(err))
	}
}
