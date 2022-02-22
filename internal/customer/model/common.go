package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/lib/pq"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"

	"time"
)

var EmptyBaseField = NewBaseField(&Subject{ID: "", Role: "", FullName: ""})

type BaseField struct {
	CreatedAt  time.Time       `db:"createdAt"`
	UpdatedAt  time.Time       `db:"updatedAt"`
	ModifiedBy *Subject        `db:"modifiedBy"`
	Version    int64           `db:"version"`
	Metadata   json.RawMessage `db:"metadata"`
}

type NullBaseField struct {
	CreatedAt  pq.NullTime      `db:"createdAt"`
	UpdatedAt  pq.NullTime      `db:"updatedAt"`
	ModifiedBy *Subject         `db:"modifiedBy"`
	Version    sql.NullInt64    `db:"version"`
	Metadata   *json.RawMessage `db:"metadata"`
}

type Subject struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"fullName"`
}

func NewBaseField(modifiedBy *Subject) BaseField {
	t := time.Now()
	return BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: modifiedBy,
		Version:    1,
		Metadata:   nsql.EmptyObjectJSON,
	}
}

type EmptyModifier = Subject

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

func ModifierModelToDTO(model Modifier) dto.Modifier {
	return dto.Modifier{
		ID:       model.ID,
		Role:     model.Role,
		FullName: model.FullName,
	}
}

func ModifierDTOToModel(dto dto.Modifier) Modifier {
	return Modifier{
		ID:       dto.ID,
		Role:     dto.Role,
		FullName: dto.FullName,
	}
}

func ItemMetadataModelToResponse(model ItemMetadata) dto.ItemMetadataResponse {
	return dto.ItemMetadataResponse{
		UpdatedAt:  model.UpdatedAt.Unix(),
		CreatedAt:  model.CreatedAt.Unix(),
		ModifiedBy: ModifierModelToDTO(*model.ModifiedBy),
		Version:    model.Version,
	}
}

func ModifierNullTime(f sql.NullTime) sql.NullTime {
	return sql.NullTime{
		Time:  f.Time,
		Valid: f.Valid,
	}
}

func ModifierNullString(f sql.NullString) sql.NullString {
	return sql.NullString{
		String: f.String,
		Valid:  f.Valid,
	}
}
