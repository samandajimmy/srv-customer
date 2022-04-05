package model

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"time"
)

type AuditLogin struct {
	BaseField
	ID           int64  `db:"id"`
	CustomerID   int64  `db:"customerId"`
	ChannelID    string `db:"channelId"`
	DeviceID     string `db:"deviceId"`
	IP           string `db:"ip"`
	Latitude     string `db:"latitude"`
	Longitude    string `db:"longitude"`
	Timestamp    string `db:"timestamp"`
	Timezone     string `db:"timezone"`
	Brand        string `db:"brand"`
	OsVersion    string `db:"osVersion"`
	Browser      string `db:"browser"`
	UseBiometric int64  `db:"useBiometric"`
	Status       int64  `db:"status"`
}

func NewAuditLogin(m *Customer, t time.Time, payload dto.LoginPayload, channelID string) AuditLogin {
	return AuditLogin{
		CustomerID:   m.ID,
		ChannelID:    channelID,
		DeviceID:     payload.DeviceID,
		IP:           payload.IP,
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
		Timestamp:    t.Format(time.RFC3339),
		Timezone:     payload.Timezone,
		Brand:        payload.Brand,
		OsVersion:    payload.OsVersion,
		Browser:      payload.Browser,
		UseBiometric: payload.UseBiometric,
		Status:       1,
		BaseField:    EmptyBaseField,
	}
}
