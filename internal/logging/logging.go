package logging

import (
	"log"
	"os"
)

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

const logLevel = TRACE

var (
	debugLogger = log.New(os.Stdout, "[DEBUG] ", 0)
	infoLogger  = log.New(os.Stdout, "[INFO]  ", 0)
	warnLogger  = log.New(os.Stdout, "[WARN]  ", 0)
	traceLogger = log.New(os.Stdout, "[TRACE] ", 0)
	errorLogger = log.New(os.Stderr, "[ERROR] ", 0)
)

func Trace(fmt string, msg ...any) {
	if logLevel <= TRACE {
		traceLogger.Printf(fmt, msg...)
	}
}

func Debug(fmt string, msg ...any) {
	if logLevel <= DEBUG {
		debugLogger.Printf(fmt, msg...)
	}
}

func Info(fmt string, msg ...any) {
	if logLevel <= INFO {
		infoLogger.Printf(fmt, msg...)
	}
}

func Warn(fmt string, msg ...any) {
	warnLogger.Printf(fmt, msg...)
}

func Error(fmt string, msg ...any) {
	errorLogger.Printf(fmt, msg...)
}
