package logger

import (
	"github.com/goburrow/gol"
)

type Logger interface {
	Trace(string)
	TraceEnabled() bool
	Debug(string)
	DebugEnabled() bool
	Info(string)
	InfoEnabled() bool
	Warn(string)
	WarnEnabled() bool
	Error(string)
	ErrorEnabled() bool
	SetLogLevel(Level)
}
type Level int

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
)

type FlogoLogger struct {
	loggerName string
	loggerImpl gol.Logger
}

func (logger *FlogoLogger) Trace(message string) {
	logger.loggerImpl.Tracef(message)
}

// TraceEnabled checks if Trace level is enabled.
func (logger *FlogoLogger) TraceEnabled() bool {
	return logger.loggerImpl.TraceEnabled()
}

// Debugf logs message at Debug level.
func (logger *FlogoLogger) Debug(message string) {
	logger.loggerImpl.Debugf(message)
}

// DebugEnabled checks if Debug level is enabled.
func (logger *FlogoLogger) DebugEnabled() bool {
	return logger.loggerImpl.DebugEnabled()
}

// Infof logs message at Info level.
func (logger *FlogoLogger) Info(message string) {
	logger.loggerImpl.Infof(message)
}

// InfoEnabled checks if Info level is enabled.
func (logger *FlogoLogger) InfoEnabled() bool {
	return logger.loggerImpl.InfoEnabled()
}

// Warnf logs message at Warning level.
func (logger *FlogoLogger) Warn(message string) {
	logger.loggerImpl.Warnf(message)
}

// WarnEnabled checks if Warning level is enabled.
func (logger *FlogoLogger) WarnEnabled() bool {
	return logger.loggerImpl.WarnEnabled()
}

// Errorf logs message at Error level.
func (logger *FlogoLogger) Error(message string) {
	logger.loggerImpl.Errorf(message)
}

// ErrorEnabled checks if Error level is enabled.
func (logger *FlogoLogger) ErrorEnabled() bool {
	return logger.loggerImpl.ErrorEnabled()
}

//SetLog Level
func (logger *FlogoLogger) SetLogLevel(logLevel Level) {
	switch logLevel {
	case Trace:
		logger.loggerImpl.(*gol.DefaultLogger).SetLevel(gol.Trace)
	case Debug:
		logger.loggerImpl.(*gol.DefaultLogger).SetLevel(gol.Debug)
	case Info:
		logger.loggerImpl.(*gol.DefaultLogger).SetLevel(gol.Info)
	case Error:
		logger.loggerImpl.(*gol.DefaultLogger).SetLevel(gol.Error)
	case Warn:
		logger.loggerImpl.(*gol.DefaultLogger).SetLevel(gol.Warn)
	default:
		logger.loggerImpl.(*gol.DefaultLogger).SetLevel(gol.Error)
	}
}

func GetLogger(name string) Logger {
	logger := gol.GetLogger(name)
	return &FlogoLogger{
		loggerName: name,
		loggerImpl: logger,
	}
}
