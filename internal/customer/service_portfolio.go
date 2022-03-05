package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger"
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
		return nil, errx.Trace(err)
	}

	return data, nil
}
