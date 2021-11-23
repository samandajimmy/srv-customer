package model

import (
	"database/sql/driver"
	"encoding/json"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
	"time"
)

type Modifier struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

func (m *Modifier) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *Modifier) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type ItemMetadata struct {
	CreatedAt  time.Time `db:"createdAt"`
	UpdatedAt  time.Time `db:"updatedAt"`
	ModifiedBy *Modifier `db:"modifiedBy"`
	Version    int64     `db:"version"`
}

func (m ItemMetadata) Upgrade(modifiedBy Modifier, opts ...time.Time) ItemMetadata {
	var t time.Time
	if len(opts) > 0 {
		t = opts[0]
	} else {
		t = time.Now()
	}

	return ItemMetadata{
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  t,
		ModifiedBy: &modifiedBy,
		Version:    m.Version + 1,
	}
}

func NewItemMetadata(modifiedBy Modifier) ItemMetadata {
	// Init timestamp
	t := time.Now()

	return ItemMetadata{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: &modifiedBy,
		Version:    1,
	}
}
