package routes

import (
	"context"
	"net/http"

	"github.com/egosha7/goph-keeper/internal/compress"
	"github.com/egosha7/goph-keeper/internal/config"
	"github.com/egosha7/goph-keeper/internal/storage"
	"github.com/egosha7/goph-keeper/server/handlers"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// SetupRoutes настраивает и возвращает обработчик HTTP-маршрутов.
func SetupRoutes(cfg *config.Config, conn *pgx.Conn, logger *zap.Logger) http.Handler {
	// Парсинг конфигурации для пула подключений
	config, err := pgxpool.ParseConfig(cfg.DataBase)
	if err != nil {
		logger.Error("Error parse config", zap.Error(err))
	}

	// Установка максимального количества соединений в пуле
	config.MaxConns = 1000

	// Создание пула подключений
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error("Error connect config", zap.Error(err))
	}

	// Создание хранилища
	repo := storage.NewPostgresURLRepository(conn, logger, pool)

	// Создание таблицы, если конфигурация БД предоставлена
	if cfg.DataBase != "" {
		repo.CreateTable()
	}

	// Создание роутера
	r := chi.NewRouter()

	// Middleware для сжатия ответа
	gzipMiddleware := compress.GzipMiddleware{}

	// Группа роутов
	r.Group(func(route chi.Router) {
		route.Use(gzipMiddleware.Apply)

		// Регистрация обработчиков для различных маршрутов
		route.Delete("/", func(w http.ResponseWriter, r *http.Request) {})
		route.Post("/auth", func(w http.ResponseWriter, r *http.Request) {
			handlers.AuthUser(w, r, repo)
		})
		route.Post("/auth/registration", func(w http.ResponseWriter, r *http.Request) {
			handlers.RegisterUser(w, r, logger, repo)
		})
		route.Post("/pass/namelist", func(w http.ResponseWriter, r *http.Request) {
			handlers.GetPasswordNameList(w, r, logger, repo)
		})
		route.Post("/card/namelist", func(w http.ResponseWriter, r *http.Request) {
			handlers.GetCardList(w, r, logger, repo)
		})
		route.Post("/pincheck", func(w http.ResponseWriter, r *http.Request) {
			handlers.CheckPinCodeHandler(w, r, logger, repo)
		})
		route.Post("/password/get", func(w http.ResponseWriter, r *http.Request) {
			handlers.GetPasswordHandler(w, r, logger, repo)
		})
		route.Post("/password/add", func(w http.ResponseWriter, r *http.Request) {
			handlers.NewPassword(w, r, logger, repo)
		})
		route.Post("/card/get", func(w http.ResponseWriter, r *http.Request) {
			handlers.GetCardHandler(w, r, logger, repo)
		})
		route.Post("/card/add", func(w http.ResponseWriter, r *http.Request) {
			handlers.NewCard(w, r, logger, repo)
		})
	})

	return r
}
