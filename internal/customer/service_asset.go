package customer

import (
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func (s *Service) AssetUploadFile(req dto.UploadRequest) (*dto.UploadResponse, error) {
	// Determine dir
	dir, err := s.AssetDirectory(req.AssetType)
	if err != nil {
		return nil, err
	}

	// generator id
	id := xid.New().String()

	// Generate filename
	fileName := req.File.Rename(id)

	// Set destination
	dest := dir + fileName

	// Upload file
	err = s.minio.Upload(req.File.File, req.File.MimeType, dest)
	if err != nil {
		s.log.Error("unable to upload file", logOption.Error(err))
		return nil, err
	}

	// Resolve file url
	fileURL := s.buildURL(dest)

	// Compose response
	resp := dto.UploadResponse{
		FileName: fileName,
		FileURL:  fileURL,
		MimeType: req.File.MimeType,
		FileSize: req.File.Header.Size,
	}

	return &resp, nil
}

func (s *Service) AssetRemoveFile(fileName string, assetType constant.AssetType) error {
	// Get directory based-on asset type
	directory, err := s.AssetDirectory(assetType)
	if err != nil {
		return errx.Trace(err)
	}

	// Remove object
	err = s.minio.Remove(directory + fileName)
	if err != nil {
		s.log.Error("error found when removing object", logOption.Error(err), logOption.Context(s.ctx))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) AssetUploadRule(assetType int) (*nhttp.UploadRule, error) {
	var rule nhttp.UploadRule

	switch assetType {
	case constant.AssetAvatarProfile, constant.AssetKTP, constant.AssetNPWP:
		rule = nhttp.UploadRule{
			Key:       nhttp.DefaultKeyFile,
			MaxSize:   nhttp.MaxFileSizeImage,
			MimeTypes: nhttp.MimeTypesImage,
		}
	default:
		return nil, constant.UnknownAssetTypeError.Trace()
	}

	return &rule, nil
}

func (s *Service) AssetGetPublicURL(assetType constant.AssetType, fileName string) string {
	// If file name is empty, return empty
	if fileName == "" {
		return ""
	}

	// Determine sub dir
	dir, err := s.AssetDirectory(assetType)
	if err != nil {
		return ""
	}

	// Set file path
	filePath := dir + fileName

	return s.buildURL(filePath)
}

func (s *Service) buildURL(filePath string) string {
	return s.config.AssetURL + "/" + filePath
}

func (s *Service) AssetDirectory(assetType int) (string, error) {
	dir, ok := constant.AssetDirs[assetType]
	if !ok {
		s.log.Errorf("unknown asset type: %d", assetType)
		return "", constant.UnknownAssetTypeError
	}
	dir += "/"
	return dir, nil
}
