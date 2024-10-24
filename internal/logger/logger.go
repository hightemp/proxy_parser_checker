package logger

import (
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

var loggers map[LogLevel]*log.Logger

func InitLogger() {
	loggers = make(map[LogLevel]*log.Logger)
	loggers[DEBUG] = log.New(os.Stdout, "DEBUG: ", log.Ltime)
	loggers[INFO] = log.New(os.Stdout, "INFO: ", log.Ltime)
	loggers[WARNING] = log.New(os.Stdout, "WARNING: ", log.Ltime)
	loggers[ERROR] = log.New(os.Stderr, "ERROR: ", log.Ltime|log.Lshortfile)
}

func LogDebug(format string, v ...any) {
	loggers[DEBUG].Printf(format, v...)
}

func LogInfo(format string, v ...any) {
	loggers[INFO].Printf(format, v...)
}

func LogWarning(format string, v ...any) {
	loggers[WARNING].Printf(format, v...)
}

func LogError(format string, v ...any) {
	loggers[ERROR].Printf(format, v...)
}

func PanicError(format string, v ...any) {
	loggers[ERROR].Panicf(format, v...)
}
