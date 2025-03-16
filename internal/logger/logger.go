package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger — глобальный логгер для проекта
var Logger = logrus.New()

func Init() {
	// Настраиваем формат логов
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Записываем логи в stdout
	Logger.SetOutput(os.Stdout)

	// Уровень логирования (можно менять через конфиг)
	Logger.SetLevel(logrus.InfoLevel)
}
