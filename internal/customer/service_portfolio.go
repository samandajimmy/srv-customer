package customer

import (
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
)

// Endpoint POST /portofolio/pds/tabemas
func (s *Service) portfolioGoldSaving(cif string) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"cif": cif,
	}

	sp := PostDataPayload{
		Path: "/portofolio/pds/tabemas",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get gold savings", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return data, nil
}
