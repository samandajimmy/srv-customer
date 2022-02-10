package customer

import (
	"context"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nredis"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
)

type Service struct {
	config       *Config
	ctx          context.Context
	repo         *RepositoryContext
	repoExternal *RepositoryContext
	log          nlogger.Logger
	responses    *ncore.ResponseMap
	redis        *nredis.Redis
}

func (h Handler) NewService(ctx context.Context) *Service {
	return &Service{
		config:       h.config,
		responses:    h.Responses,
		redis:        h.redis,
		repo:         h.repo.WithContext(ctx),
		repoExternal: h.repoExternal.WithContext(ctx),
		ctx:          ctx,
		log:          nlogger.Get().NewChild(nlogger.Context(ctx)),
	}
}

func (s *Service) GetOrderBy(sortBy string, sortDirection string, rules []string) (string, string) {
	if nval.InArrayString(sortBy, rules) {
		// Normalize direction
		sortDirection = strings.ToUpper(sortDirection)
		if sd := sortDirection; sd != `ASC` && sd != `DESC` {
			sortDirection = `ASC`
		}
		return sortBy, sortDirection
	}

	return `createdAt`, `DESC`
}

func (s *Service) Close() {
	// Close database connection to free pool
	err := s.repo.conn.Close()
	if err != nil {
		s.log.Error("Failed to close connection", nlogger.Error(err))
	}
	err = s.repoExternal.conn.Close()
	if err != nil {
		s.log.Error("Failed to close connection", nlogger.Error(err))
	}

}
