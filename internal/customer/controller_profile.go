package customer

import (
	"errors"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/validate"
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
	userRefID := GetUserRefID(rx)

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
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.UpdateProfilePayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PutUpdateProfile(&payload)
	if err != nil {
		log.Error("Bad request validate payload")
		return nil, err
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

func (c *ProfileController) PostUpdateAvatar(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get multipart file
	file, err := nhttp.GetFile(rx.Request, constant.KeyUserFile, nhttp.MaxFileSizeImage, nhttp.MimeTypesImage)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, constant.NoFileError.Trace()
		}
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Set payload avatar
	payloadUserFile := dto.UploadUserFilePayload{
		File:      file,
		AssetType: constant.AssetAvatarProfile,
	}

	// Upload Userfile
	userFile, err := svc.UploadUserFile(payloadUserFile)
	if err != nil {
		log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Set payload update avatar
	updateAvatar := dto.UpdateAvatarPayload{
		UpdateUserFile: dto.UpdateUserFile{
			FileName:  userFile.FileName,
			UserRefID: userRefID,
			AssetType: payloadUserFile.AssetType,
		},
		FileSize: userFile.FileSize,
		MimeType: userFile.MimeType,
	}

	// Call service
	err = svc.UpdateAvatar(updateAvatar)
	if err != nil {
		log.Error("error when call update avatar service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(userFile), nil
}

func (c *ProfileController) PostUpdateKTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get multipart file
	file, err := nhttp.GetFile(rx.Request, constant.KeyUserFile, nhttp.MaxFileSizeImage, nhttp.MimeTypesImage)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, constant.NoFileError.Trace()
		}
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Set payload avatar
	payloadUserFile := dto.UploadUserFilePayload{
		File:      file,
		AssetType: constant.AssetKTP,
	}

	// Upload Userfile
	userFile, err := svc.UploadUserFile(payloadUserFile)
	if err != nil {
		log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Call service
	err = svc.UpdateIdentity(dto.UpdateUserFile{
		FileName:  userFile.FileName,
		UserRefID: userRefID,
		AssetType: constant.AssetKTP,
	})
	if err != nil {
		log.Error("error when call update avatar service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(userFile), nil
}

func (c *ProfileController) PostUpdateNPWP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get multipart file
	file, err := nhttp.GetFile(rx.Request, constant.KeyUserFile, nhttp.MaxFileSizeImage, nhttp.MimeTypesImage)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, constant.NoFileError.Trace()
		}
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Set payload for upload userfile
	payloadUserFile := dto.UploadUserFilePayload{
		File:      file,
		AssetType: constant.AssetNPWP,
	}

	// Upload userfile
	userFile, err := svc.UploadUserFile(payloadUserFile)
	if err != nil {
		log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Get payload
	payload := dto.UpdateNPWPPayload{
		NoNPWP:    rx.FormValue("no_npwp"),
		UserRefID: userRefID,
		FileName:  userFile.FileName,
	}

	// Validate payload
	err = validate.PostUpdateNPWP(&payload)
	if err != nil {
		log.Error("Bad request validate payload")
		return nil, err
	}

	// Call service
	err = svc.UpdateNPWP(payload)
	if err != nil {
		log.Error("error when call update npwp service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(userFile), nil
}

func (c *ProfileController) PostUpdateSID(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get multipart file
	file, err := nhttp.GetFile(rx.Request, constant.KeyUserFile, nhttp.MaxFileSizeImage, nhttp.MimeTypesImage)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, constant.NoFileError.Trace()
		}
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Set payload SID
	payloadUserFile := dto.UploadUserFilePayload{
		File:      file,
		AssetType: constant.AssetNPWP,
	}
	// Upload Userfile
	userFile, err := svc.UploadUserFile(payloadUserFile)
	if err != nil {
		log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Get payload
	payload := dto.UpdateSIDPayload{
		NoSID:     rx.FormValue("no_sid"),
		UserRefID: userRefID,
		FileName:  userFile.FileName,
	}

	// Validate payload
	err = validate.PostUpdateSID(&payload)
	if err != nil {
		log.Error("Bad request validate payload")
		return nil, err
	}

	// Call service
	err = svc.UpdateSID(payload)
	if err != nil {
		log.Error("error when call update SID service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(userFile), nil
}

func (c *ProfileController) GetStatus(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

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
