package customer

import (
	"crypto/md5"
	"fmt"
	"github.com/rs/xid"
	"strings"
	"time"
)

func generateReferralCode(prefix string) string {
	unique := xid.New().String()
	rand := fmt.Sprintf("%s", unique[0:7])

	var referralCode string
	if prefix == "PDS" {
		referralCode = prefix + rand
	} else {
		referralCode = prefix
	}

	return strings.ToUpper(referralCode)

}

func monthsToUnix(month int) int64 {
	twoMonth := time.Now().AddDate(0, month, 0).Unix()
	return twoMonth
}

func hoursToSeconds(hour int64) int64 {
	now := time.Now()
	return now.Add(time.Hour * time.Duration(hour)).Unix()
}

func stringToMD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
