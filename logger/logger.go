package logger

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

func GetLoggerFactory() LoggerFactory {
	return logFactory
}
