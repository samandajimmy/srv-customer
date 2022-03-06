package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type ProfileController struct {
	*Handler
}

func NewProfileController(h *Handler) *ProfileController {
	return &ProfileController{
		h,
	}
}

func (c *ProfileController) GetDetail(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.CustomerProfile(userRefID)
	if err != nil {
		log.Error("error when call customer profile service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *ProfileController) PutUpdate(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get payload
	var payload dto.UpdateProfileRequest
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	err = svc.UpdateCustomerProfile(userRefID, payload)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage("Update data user berhasil").SetData(false), nil
}

func (c *ProfileController) UpdateAvatar(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.UpdateAvatar(dto.UpdateUserFile{
		Request:   rx.Request,
		UserRefID: userRefID,
		AssetType: constant.AssetAvatarProfile,
	})
	if err != nil {
		log.Error("error when call update avatar service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *ProfileController) UpdateKTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Error("error user auth", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()
	// Call service
	resp, err := svc.UpdateIdentity(dto.UpdateUserFile{
		Request:   rx.Request,
		UserRefID: userRefID,
		AssetType: constant.AssetKTP,
	})
	if err != nil {
		log.Error("error when call update avatar service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *ProfileController) UpdateNPWP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Error("error user auth", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get payload
	payload := dto.UpdateNPWPRequest{
		Request:   rx.Request,
		NoNPWP:    rx.FormValue("no_npwp"),
		UserRefID: userRefID,
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()
	// Call service
	resp, err := svc.UpdateNPWP(payload)
	if err != nil {
		log.Error("error when call update npwp service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *ProfileController) UpdateSID(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()
	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Error("error user auth", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}
	// Get payload
	var payload dto.UpdateSIDRequest
	number := rx.FormValue("no_sid")
	// Validate payload
	payload.NoSID = number
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}
	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()
	// Call service
	resp, err := svc.UpdateSID(dto.UpdateSIDRequest{
		Request:   rx.Request,
		NoSID:     number,
		UserRefID: userRefID,
	})
	if err != nil {
		log.Error("error when call update SID service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *ProfileController) CheckStatus(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.CheckStatus(userRefID)
	if err != nil {
		log.Error("error found when call check status service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
