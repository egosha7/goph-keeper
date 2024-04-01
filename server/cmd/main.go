package main

import (
	"context"
	"fmt"
	"github.com/egosha7/goph-keeper/internal/config"
	"github.com/egosha7/goph-keeper/internal/db"
	loger "github.com/egosha7/goph-keeper/internal/logger"
	"github.com/egosha7/goph-keeper/internal/router"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Глобальные переменные
var (
	// Version - это версия сборки приложения.
	Version string
	// BuildTime - это временная метка времени сборки приложения.
	BuildTime string
	// Commit - это хеш коммита приложения.
	Commit string
)

// main - это основная точка входа для службы shortlink.
func main() {
	fmt.Printf("Версия сборки: %s\n", Version)
	fmt.Printf("Дата сборки: %s\n", BuildTime)
	fmt.Printf("Коммит: %s\n", Commit)

	// Настройка логгера.
	logger, err := loger.SetupLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания логгера: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Проверка конфигурации из флагов и переменных окружения.
	cfg := config.OnFlag(logger)

	// Подключение к базе данных.
	conn, err := db.ConnectToDB(cfg)
	if err != nil {
		logger.Error("Ошибка подключения к базе данных", zap.Error(err))
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Настройка маршрутов для приложения.
	r := routes.SetupRoutes(cfg, conn, logger)

	// Настройка обработки сигналов для грациозного завершения.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var wg sync.WaitGroup

	// Запуск горутины для обработки сигналов.
	go func() {
		sig := <-signalCh
		fmt.Printf("Получен сигнал %v. Завершение работы...\n", sig)

		// Дождемся завершения оставшихся запросов.
		wg.Wait()

		// Завершаем программу.
		os.Exit(0)
	}()

	if err := http.ListenAndServe(cfg.Addr, loger.LogMiddleware(logger, r)); err != nil {
		logger.Error("Ошибка запуска HTTP сервера", zap.Error(err))
		os.Exit(1)
	}
}
