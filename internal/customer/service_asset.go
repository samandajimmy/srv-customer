package customer

import (
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
		s.log.Error("unable to upload file", err)
		return nil, err
	}

	// Resolve file url
	fileUrl := s.buildUrl(dest)

	// Compose response
	resp := dto.UploadResponse{
		FileName: fileName,
		FileUrl:  fileUrl,
	}

	return &resp, nil
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
		return nil, s.responses.GetError("E_AST_1")
	}

	return &rule, nil
}

func (s *Service) AssetGetPublicUrl(assetType int, fileName string) string {
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

	return s.buildUrl(filePath)
}

func (s *Service) buildUrl(filePath string) string {
	return s.config.AssetUrl + "/" + filePath
}

func (s *Service) AssetDirectory(assetType int) (string, error) {
	dir, ok := constant.AssetDirs[assetType]
	if !ok {
		err := s.responses.GetError("E_AST_1")
		s.log.Errorf("unknown asset type: %d", assetType)
		return "", err
	}
	dir += "/"
	return dir, nil
}
