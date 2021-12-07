package model

import (
	"encoding/json"
)

type AuditLogin struct {
	Id           int64           `db:"id"`
	CustomerId   int64           `db:"customerId"`
	ChannelId    string          `db:"channelId"`
	DeviceId     string          `db:"deviceId"`
	IP           string          `db:"ip"`
	Latitude     string          `db:"latitude"`
	Longitude    string          `db:"longitude"`
	Timestamp    string          `db:"timestamp"`
	Timezone     string          `db:"timezone"`
	Brand        string          `db:"brand"`
	OsVersion    string          `db:"osVersion"`
	Browser      string          `db:"browser"`
	UseBiometric int64           `db:"useBiometric"`
	Status       int64           `db:"status"`
	Metadata     json.RawMessage `db:"metadata"`
	ItemMetadata
}
