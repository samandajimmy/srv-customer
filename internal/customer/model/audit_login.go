package model

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
