package main

import (
	"go.uber.org/zap"
)

// NewConsoleLogger Разместил здесь, потому что иначе тесты ругаются
func NewConsoleLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic("cannot initialize zap")
	}
	// это нужно добавить, если логер буферизован
	// в данном случае не буферизован, но привычка хорошая
	//не знаю, как обрабатывать ошибку
	defer func() {
		_ = logger.Sync()
	}()

	return logger
}
