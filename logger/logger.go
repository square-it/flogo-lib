package logger

import "errors"

type Logger interface {
	Debug(...interface{})
	DebugEnabled() bool
	Info(...interface{})
	InfoEnabled() bool
	Warn(...interface{})
	WarnEnabled() bool
	Error(...interface{})
	ErrorEnabled() bool
	SetLogLevel(Level)
}

type LoggerFactory interface {
	GetLogger(name string) (Logger, error)
}

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

var logFactory LoggerFactory

func RegisterLoggerFactory(factory LoggerFactory) {
	logFactory = factory
}

func GetLogger(name string) (Logger, error) {
	if logFactory == nil {
		return nil, errors.New("No logger factory found.")
	}
	return logFactory.GetLogger(name)
}
