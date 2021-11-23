package nlogger

import (
	"errors"
	stdLog "log"
	"os"
	"testing"
)

// Logger instance
var testLogger Logger

// Init sample variables
var metadata = map[string]interface{}{
	"string":  "string",
	"integer": 0,
	"boolean": true,
	"array":   []int{1, 2, 3, 4, 5},
	"struct": struct {
		Text    string
		Boolean bool
	}{
		Text:    "text",
		Boolean: false,
	},
}

func TestMain(m *testing.M) {
	testLogger = NewStdLogger(LevelDebug, os.Stdout, "", stdLog.LstdFlags)

	// Run Test
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

func TestFatal(t *testing.T) {
	testLogger.Fatal("Testing FATAL with message only")
	testLogger.Fatalf("Testing FATAL with formatted message: %s %s", "arg1", "arg2")
	testLogger.Fatal("Testing FATAL with options. Formatted Message: %s %s %s", Options{
		Error:    errors.New("a fatal error occurred"),
		Metadata: metadata,
		FmtArgs:  []interface{}{"arg1", "arg2", "arg3"},
	})
}

func TestError(t *testing.T) {
	testLogger.Error("Testing ERROR with message only")
	testLogger.Errorf("Testing ERROR with formatted message: %s %s", "arg1", "arg2")
	testLogger.Error("Testing ERROR with options. Formatted Message: %s %s %s", Options{
		Error:    errors.New("an error occurred"),
		Metadata: metadata,
		FmtArgs:  []interface{}{"arg1", "arg2", "arg3"},
	})
}

func TestWarn(t *testing.T) {
	testLogger.Warn("Testing WARN with message only")
	testLogger.Warnf("Testing WARN with formatted message: %s %s", "arg1", "arg2")
	testLogger.Warn("Testing WARN with options. Formatted Message: %s %s %s", Options{
		Metadata: metadata,
		FmtArgs:  []interface{}{"arg1", "arg2", "arg3"},
	})
}

func TestInfo(t *testing.T) {
	testLogger.Info("Testing INFO with message only")
	testLogger.Infof("Testing INFO with formatted message: %s %s", "arg1", "arg2")
	testLogger.Info("Testing INFO with options. Formatted Message: %s %s %s", Options{
		Metadata: metadata,
		FmtArgs:  []interface{}{"arg1", "arg2", "arg3"},
	})
}

func TestDebug(t *testing.T) {
	testLogger.Debug("Testing DEBUG with message only")
	testLogger.Debugf("Testing DEBUG with formatted message: %s %s", "arg1", "arg2")
	testLogger.Debug("Testing DEBUG with options. Formatted Message: %s %s %s", Options{
		Metadata: metadata,
		FmtArgs:  []interface{}{"arg1", "arg2", "arg3"},
	})
}

func TestDefault(t *testing.T) {
	l := Get()
	l.Error("This is called from default logger")
}
