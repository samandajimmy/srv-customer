package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
)

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
