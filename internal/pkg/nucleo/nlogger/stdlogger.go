package nlogger

import (
	"encoding/json"
	"errors"
	"io"
	stdLog "log"
	"os"
)

type StdLogger struct {
	level       int
	levelPrefix map[int]string
	skipTrace   int
	writer      *stdLog.Logger
}

func (l StdLogger) Fatal(msg string, args ...interface{}) {
	l.print(LevelFatal, msg, GetOptions(args))
}

func (l StdLogger) Fatalf(format string, args ...interface{}) {
	l.print(LevelFatal, format, &Options{
		FmtArgs: args,
	})
}

func (l StdLogger) Error(msg string, args ...interface{}) {
	l.print(LevelError, msg, GetOptions(args))
}

func (l StdLogger) Errorf(format string, args ...interface{}) {
	l.print(LevelError, format, &Options{
		FmtArgs: args,
	})
}

func (l StdLogger) Warn(msg string, args ...interface{}) {
	l.print(LevelWarn, msg, GetOptions(args))
}

func (l StdLogger) Warnf(format string, args ...interface{}) {
	l.print(LevelWarn, format, &Options{
		FmtArgs: args,
	})
}

func (l StdLogger) Info(msg string, args ...interface{}) {
	l.print(LevelInfo, msg, GetOptions(args))
}

func (l StdLogger) Infof(format string, args ...interface{}) {
	l.print(LevelInfo, format, &Options{
		FmtArgs: args,
	})
}

func (l *StdLogger) Debug(msg string, args ...interface{}) {
	l.print(LevelDebug, msg, GetOptions(args))
}

func (l *StdLogger) Debugf(format string, args ...interface{}) {
	l.print(LevelDebug, format, &Options{
		FmtArgs: args,
	})
}

func NewStdLogger(level int, w io.Writer, prefix string, flags int) Logger {
	// If writer is nil, set default writer to Stdout
	if w == nil {
		w = os.Stdout
	}

	// Init standard logger instance
	l := StdLogger{
		level: level,
		levelPrefix: map[int]string{
			LevelPanic: "[PANIC] ",
			LevelFatal: "[FATAL] ",
			LevelError: "[ERROR] ",
			LevelWarn:  " [WARN] ",
			LevelInfo:  " [INFO] ",
			LevelDebug: "[DEBUG] ",
		},
		skipTrace: 2,
		writer:    stdLog.New(w, prefix, flags),
	}
	return &l
}

func (l *StdLogger) print(outLevel int, msg string, options *Options) {
	// if output level is greater than log level, don't print
	if outLevel > l.level {
		return
	}

	// Generate prefix
	prefix := l.levelPrefix[outLevel]

	// If options is exist
	if options != nil {
		// If formatted arguments is available, then print as formatted
		if options.FmtArgs != nil && len(options.FmtArgs) > 0 {
			l.writer.Printf(prefix+msg+"\n", options.FmtArgs...)
		} else {
			l.writer.Printf("%s%s\n", prefix, msg)
		}

		// If error exists, then print error and its trace
		if options.Error != nil && outLevel <= LevelError {
			filePath, line := Trace(l.skipTrace)
			l.writer.Printf("  > Error: %s\n", options.Error)
			l.writer.Printf("  > Trace: %s:%d\n", filePath, line)
			// Print cause
			if unErr := errors.Unwrap(options.Error); unErr != nil {
				l.writer.Printf("  > ErrorCause: %s\n", unErr)
			}
		}

		if options.Metadata != nil && len(options.Metadata) > 0 {
			// Serialize to json
			metadata, err := json.MarshalIndent(options.Metadata, "  ", "  ")
			// If not error, then print
			if err == nil {
				l.writer.Printf("  > Metadata: \n  %s\n", metadata)
			}
		}

		return
	}

	l.writer.Printf("%s%s\n", prefix, msg)
}
