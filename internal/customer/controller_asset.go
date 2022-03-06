package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type Asset struct {
	*Handler
}

func NewAsset(h *Handler) *Asset {
	h.Config.AssetURL = h.Config.MinioURL + "/" + h.Config.MinioBucket

	return &Asset{h}
}

func (h *Asset) UploadFile(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get query
	q := rx.URL.Query()

	// Init service
	svc := h.NewService(ctx)
	defer svc.Close()

	// Get asset type
	assetType := nval.ParseIntFallback(q.Get("asset_type"), 0)

	// Get rules by asset type
	rule, err := svc.AssetUploadRule(assetType)
	if err != nil {
		log.Error("error found when get rules by asset type", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Parse multipart file
	file, err := nhttp.GetFile(rx.Request, rule.Key, rule.MaxSize, rule.MimeTypes)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Upload file payload
	filePayload := dto.UploadRequest{
		AssetType: assetType,
		File:      file,
	}

	// Upload a file
	resp, err := svc.AssetUploadFile(filePayload)
	if err != nil {
		log.Error("error found when call service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(resp), nil
}
