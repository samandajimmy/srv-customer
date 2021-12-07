package ntime

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"time"
)

func NewTimeWIB(t time.Time) time.Time {
	loc, err := time.LoadLocation(constant.WIB)
	if err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}
	return t
}

func ChangeToUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}

func ChangeTimezone(t time.Time, timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}
	return t
}
