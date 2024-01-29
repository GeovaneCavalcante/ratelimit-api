package logger

import "log/slog"

func Info(message string) {
	logger := slog.Default()
	logger.Info(message)
}

func Error(message string, err error) {
	logger := slog.Default()
	if err == nil {
		logger.Error(message)
		return
	}
	logger.Error(message + " : " + err.Error())
}

func Warn(message string, err error) {
	logger := slog.Default()
	if err == nil {
		logger.Warn(message)
		return
	}
	logger.Warn(message + " : " + err.Error())
}
