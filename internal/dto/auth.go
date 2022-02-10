package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
)

type Subject struct {
	SubjectID    string
	SubjectRefID int64
	SubjectRole  string
	SubjectType  constant.SubjectType
	ModifiedBy   Modifier
	Metadata     map[string]string
	SessionID    int64
}

type ClientCredential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
