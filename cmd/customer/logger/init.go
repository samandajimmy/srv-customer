package logger

import (
	jsonLogger "github.com/nbs-go/nlogger-json"
	"github.com/nbs-go/nlogger/v2"
	"os"
)

func init() {
	// Register json logger
	lv := os.Getenv(nlogger.EnvLogLevel)
	nlogger.Register(jsonLogger.New("srv-customer", lv, os.Stdout))
}
