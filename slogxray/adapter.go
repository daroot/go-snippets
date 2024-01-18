// Package slogxray adapts slog logger for use in aws-xray-sdk-go.
package slogxray

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-xray-sdk-go/xraylog"
)

// Adapter converts an slog.Logger into something xraylog can accept.
type Adapter struct {
	log *slog.Logger
}

// Log outputs an xray log event via the configure zerolog.Logger.
func (a Adapter) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	s := msg.String()
	switch level {
	case xraylog.LogLevelDebug:
		a.log.Debug(s)
	case xraylog.LogLevelInfo:
		a.log.Info(s)
	case xraylog.LogLevelWarn:
		a.log.Warn(s)
	case xraylog.LogLevelError:
		a.log.Error(s)
	}
}

// New creates a new zerologxray.Adapter from a zerolog.Logger.
func New(logger *slog.Logger) *Adapter {
	return &Adapter{log: logger}
}
