package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net"
	"regexp"
)

// Config - структура конфигурации приложения
type Config struct {
	Addr     string `env:"SERVER_ADDRESS" json:"server_address"` // Адрес сервера
	BaseURL  string `env:"BASE_URL" json:"base_url"`             // Базовый адрес результирующего сокращенного URL
	DataBase string `env:"DATABASE_DSN" json:"database_dsn"`     // Адрес базы данных
}

// Default - функция для создания новой конфигурации с значениями по умолчанию
func Default() *Config {
	return &Config{
		Addr:     "localhost:8080",
		BaseURL:  "http://localhost:8080",
		DataBase: "postgres://postgres:egosha@localhost:5432/keeper", // postgres://egosha:admin@localhost:5432/keeper
	}
}

// OnFlag - функция для чтения значений из флагов командной строки и записи их в структуру Config
func OnFlag(logger *zap.Logger) *Config {
	defaultValue := Default()

	// Инициализация флагов командной строки
	config := Config{}
	flag.StringVar(&config.Addr, "a", defaultValue.Addr, "HTTP-адрес сервера")
	flag.StringVar(&config.BaseURL, "b", defaultValue.BaseURL, "Базовый адрес результирующего сокращенного URL")
	flag.StringVar(&config.DataBase, "d", defaultValue.DataBase, "Адрес базы данных")
	flag.Parse()

	godotenv.Load()

	// Парсинг переменных окружения в структуру Config
	if err := env.Parse(&config); err != nil {
		logger.Error("Ошибка при парсинге переменных окружения", zap.Error(err))
	}

	// Проверка корректности введенных значений флагов
	if _, _, err := net.SplitHostPort(config.Addr); err != nil {
		panic(err)
	}
	if matched, _ := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, config.BaseURL); !matched {
		panic("Invalid base URL")
	}

	return &config
}
