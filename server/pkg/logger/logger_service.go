package logger

import "log/slog"

type LoggerService struct {
	logger *slog.Logger
}

func NewLoggerService(logger *slog.Logger) *LoggerService {
	return &LoggerService{
		logger: logger,
	}
}

func (ls *LoggerService) Info(msg string, args ...any) {
	ls.logger.Info(msg, args...)
}

func (ls *LoggerService) Error(msg string, args ...any) {
	ls.logger.Error(msg, args...)
}
