package lsyslog

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-syslog"
	"github.com/juju/loggo"
)

var formatter = &SyslogFormatter{}

// SyslogFormatter provides a simple concatenation of message components.
type SyslogFormatter struct{}

// Format returns the parameters separated by spaces except for filename and
// line which are separated by a colon.
func (*SyslogFormatter) Format(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) string {

	// Just get the basename from the filename
	filename = filepath.Base(filename)
	return fmt.Sprintf("%s %s:%d %s", module, filename, line, message)
}

type syslogWriter struct {
	syslogger gsyslog.Syslogger
}

// NewSyslogWriter returns a new writer that writes
// log messages to syslog in a simple format tailored for syslog
func NewSyslogWriter(p gsyslog.Priority, facility, tag string) loggo.Writer {
	syslogger, err := gsyslog.NewLogger(p, facility, tag)
	if err != nil {
		panic(err)
	}
	return &syslogWriter{syslogger}
}

// NewDefaultSyslogWriter returns a new writer that writes
// log messages to syslog in a simple format tailored for syslog.
// Note this defaults to using LOCAL7.
func NewDefaultSyslogWriter(level loggo.Level, tag, facility string) loggo.Writer {
	if facility == "" {
		facility = "LOCAL7"
	}
	syslogger, err := gsyslog.NewLogger(convertLevel(level), facility, tag)
	if err != nil {
		panic(err) // of course this will never happen
	}
	return &syslogWriter{syslogger}
}

func (slog *syslogWriter) Write(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) {
	logLine := formatter.Format(level, module, filename, line, timestamp, message)
	slog.syslogger.WriteLevel(convertLevel(level), []byte(logLine))
}

func convertLevel(level loggo.Level) gsyslog.Priority {
	switch level {
	case loggo.DEBUG:
		return gsyslog.LOG_DEBUG
	case loggo.INFO:
		return gsyslog.LOG_INFO
	case loggo.WARNING:
		return gsyslog.LOG_WARNING
	case loggo.CRITICAL:
		return gsyslog.LOG_CRIT
	case loggo.ERROR:
		return gsyslog.LOG_ERR
	case loggo.TRACE:
		return gsyslog.LOG_DEBUG
	default:
		return gsyslog.LOG_DEBUG
	}

}
