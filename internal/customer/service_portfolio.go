package customer

import (
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

// Endpoint POST /portofolio/pds/tabemas
func (s *Service) portfolioGoldSaving(cif string) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"cif": cif,
	}

	sp := PostDataPayload{
		Url:  "/portofolio/pds/tabemas",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get gold savings", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, ncore.TraceError("error found when get gold savings", err)
	}

	return data, nil
}
