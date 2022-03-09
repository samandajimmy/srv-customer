package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
)

type Subject struct {
	ID          string
	RefID       int64
	Role        string
	FullName    string
	SubjectType constant.SubjectType
	SessionID   int64
	Metadata    map[string]string
}

func (s *Subject) ModifiedBy() *Modifier {
	return &Modifier{
		ID:       s.ID,
		Role:     s.Role,
		FullName: s.FullName,
	}
}

type Modifier struct {
	ID       string                `json:"-"`
	Role     constant.ModifierRole `json:"role"`
	FullName string                `json:"fullName"`
}

type ItemMetadataResponse struct {
	CreatedAt  int64    `json:"createdAt"`
	UpdatedAt  int64    `json:"updatedAt"`
	ModifiedBy Modifier `json:"modifiedBy"`
	Version    int64    `json:"version"`
}
