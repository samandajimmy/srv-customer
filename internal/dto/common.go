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

type FindOptions struct {
	Limit         int                    `json:"limit"`
	Skip          int                    `json:"skip"`
	SortBy        string                 `json:"sortBy"`
	SortDirection string                 `json:"sortDirection"`
	Filters       map[string]interface{} `json:"-"`
}

type ListMetadata struct {
	Count int64 `json:"count"`
	FindOptions
}
