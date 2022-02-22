package dto

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"

type UploadRequest struct {
	AssetType int
	File      nhttp.MultipartFile
	DestDir   string
}

type UploadResponse struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
}
