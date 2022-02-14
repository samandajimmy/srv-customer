package customer

import (
	"fmt"
	"strings"
	"time"
)

func generateReferralCode(prefix string) string {
	unique := uniqID("")
	rand := fmt.Sprintf("%s", unique[8:15])

	var referralCode string
	if prefix == "PDS" {
		referralCode = prefix + rand
	} else {
		referralCode = prefix
	}

	return strings.ToUpper(referralCode)

}

func uniqID(prefix string) string {
	now := time.Now()
	sec := now.Unix()
	usec := now.UnixNano() % 0x100000
	return fmt.Sprintf("%s%08x%05x", prefix, sec, usec)
}

func monthsToSeconds(month int) int {
	now := time.Now()
	return now.AddDate(0, month, 0).Second()
}

func hoursToSeconds(hour int64) int {
	now := time.Now()
	return now.Add(time.Hour * time.Duration(hour)).Second()
}
