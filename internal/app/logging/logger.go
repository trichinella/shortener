package logging

import (
	"go.uber.org/zap"
	"log"
	"sync"
)

var once sync.Once
var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {
	once.Do(func() {
		//в далеком будущем здесь будет расхождение на контуры, смотреть будет на переменную окружения
		Logger = newConsoleLogger()
		Sugar = Logger.Sugar()
	})
}

// newConsoleLogger
func newConsoleLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		log.Fatalf("cannot initialize zap")
	}

	return logger
}
