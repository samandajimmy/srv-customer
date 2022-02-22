package customer

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"github.com/rs/xid"
	"strings"
	"time"
)

func generateReferralCode(prefix string) string {
	unique := xid.New().String()
	rand := unique[0:6]

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

func stringToMD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str))) //nolint:gosec
}
