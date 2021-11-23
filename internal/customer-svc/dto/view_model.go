package dto

type RedirectPageResponse struct {
	Title                 string
	Domain                string
	Description           string
	ThumbnailURL          string
	OriginalURL           string
	LogoURL               string
	BaseStaticURL         string
	ConfLanguage          string
	ConfCountdownSecStart int
}

type ErrorPageData struct {
	Title         string
	ErrorCode     int
	ErrorMessage  string
	ImageURL      string
	ImageAltText  string
	BaseStaticURL string
	ConfLanguage  string
}
