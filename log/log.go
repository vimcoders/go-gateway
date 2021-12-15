package log

import (
	"github.com/vimcoders/go-driver"
)

var logger, _ = driver.NewSyslogger()

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	logger.Warning(format, v...)
}
