package nlogger

import (
	stdLog "log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Log Level constants as defined in RFC5424.
const (
	LevelPanic = iota
	LevelFatal
	_ // LevelCritical
	LevelError
	LevelWarn
	_ // LevelNotice
	LevelInfo
	LevelDebug
	_ // LevelTrace
)

// Configuration constants.
const (
	// Log config env key
	EnvLogLevel = "LOG_LEVEL"
	// Default Values
	DefaultLevel = LevelError
)

func Error(err error) Options {
	return Options{
		Metadata: make(map[string]interface{}),
		Error:    err,
	}
}

func NewOptions() Options {
	return Options{
		Metadata: make(map[string]interface{}),
	}
}

/// Options is available Options to print log
type Options struct {
	Metadata map[string]interface{}
	Error    error
	FmtArgs  []interface{}
}

func (o Options) AddMetadata(key string, value interface{}) Options {
	o.Metadata[key] = value
	return o
}

func (o Options) SetFormat(args ...interface{}) Options {
	o.FmtArgs = args
	return o
}

func (o Options) SetError(err error) Options {
	o.Error = err
	return o
}

/// Logger contract defines methods that must be available for a Logger.
///
/// Fatal must write an error, message that explaining the error and where its occurred in FATAL level.
/// To trace message, use Trace function and skip 1.
///
/// Fatalf must write a formatted message and where its occurred in FATAL level.
/// To trace message, use Trace function and skip 1.
///
/// Error must write an error, message that explaining the error and where its occurred in ERROR level.
/// To trace message, use Trace function and skip 1.
///
/// Errorf must write a formatted message and where its occurred in ERROR level.
/// To trace message, use Trace function and skip 1.
///
/// Warn must write a message in WARN level.
///
/// Warnf must write a formatted message in WARN level.
///
/// Info must write a message in INFO level.
///
/// Infof must write a formatted message in INFO level.
///
/// Debug must write a message in DEBUG level.
///
/// Debugf must write a formatted message in DEBUG level.
type Logger interface {
	Fatal(msg string, options ...interface{})
	Fatalf(format string, args ...interface{})
	Error(msg string, options ...interface{})
	Errorf(format string, args ...interface{})
	Warn(msg string, options ...interface{})
	Warnf(format string, args ...interface{})
	Info(msg string, options ...interface{})
	Infof(format string, args ...interface{})
	Debug(msg string, options ...interface{})
	Debugf(format string, args ...interface{})
}

/// log is a singleton logger instance
var log Logger
var logMutex sync.RWMutex

/// Get retrieve singleton logger instance
func Get() Logger {
	// If log is nil, initiate standard logger
	if log == nil {
		// Get logger from env
		logLevel := DefaultLevel
		logLevelStr, ok := os.LookupEnv(EnvLogLevel)
		if ok {
			logLevel = ParseLevel(logLevelStr)
		}

		// Init standard logger
		l := NewStdLogger(logLevel, os.Stderr, "", stdLog.LstdFlags)

		// Register logger
		Register(l)
		log.Debug("No logger found. StdLogger initiated")
	}
	return log
}

/// Register logger instance
func Register(l Logger) {
	// If logger is nil, return error
	if l == nil {
		panic("nbs-go/nucleo/nlogger: logger to be registered is nil")
	}

	// Set logger
	logMutex.Lock()
	defer logMutex.Unlock()
	log = l
}

/// Trace retrieve where the code is being called and returns full path of file where the error occurred
func Trace(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		file = "<???>"
		line = 1
	}
	return file, line
}

/// ParseLevel parse level from string to Log Level enum
func ParseLevel(level string) int {
	switch strings.ToLower(level) {
	case "fatal", "1":
		return LevelFatal
	case "warn", "4":
		return LevelWarn
	case "info", "6":
		return LevelInfo
	case "debug", "7":
		return LevelDebug
	default:
		return LevelError
	}
}

/// GetOptions retrieve logger.Options from arguments spread
func GetOptions(args []interface{}) *Options {
	// If args is empty, return nil
	argsLen := len(args)
	if argsLen == 0 {
		return nil
	}

	// Set 1st argument as logger options
	opts, ok := args[0].(Options)

	// If options, then return
	if ok {
		return &opts
	}

	// Else, set format args
	return &Options{
		FmtArgs: args,
	}
}
