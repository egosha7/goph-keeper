package routes

import (
	"context"
	"github.com/egosha7/goph-keeper/internal/compress"
	"github.com/egosha7/goph-keeper/internal/config"
	"github.com/egosha7/goph-keeper/internal/handlers"
	"github.com/egosha7/goph-keeper/internal/repository"
	"github.com/egosha7/goph-keeper/internal/service"
	"net/http"

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
	repo := repository.NewPostgreSQLRepository(pool, logger)
	services := service.NewUserService(repo, logger)
	h := handlers.NewHandler(services, logger)

	// Создание роутера
	r := chi.NewRouter()

	// Middleware для сжатия ответа
	gzipMiddleware := compress.GzipMiddleware{}

	// Группа роутов
	r.Group(
		func(route chi.Router) {
			route.Use(gzipMiddleware.Apply)

			// Регистрация обработчиков для различных маршрутов
			route.Delete("/", func(w http.ResponseWriter, r *http.Request) {})
			route.Post(
				"/auth", func(w http.ResponseWriter, r *http.Request) {
					h.AuthUser(w, r)
				},
			)
			route.Post(
				"/auth/registration", func(w http.ResponseWriter, r *http.Request) {
					h.RegisterUser(w, r)
				},
			)
			route.Post(
				"/pass/namelist", func(w http.ResponseWriter, r *http.Request) {
					h.GetPasswordNameList(w, r)
				},
			)
			route.Post(
				"/card/namelist", func(w http.ResponseWriter, r *http.Request) {
					h.GetCardList(w, r)
				},
			)
			route.Post(
				"/pincheck", func(w http.ResponseWriter, r *http.Request) {
					h.CheckPinCodeHandler(w, r)
				},
			)
			route.Post(
				"/password/get", func(w http.ResponseWriter, r *http.Request) {
					h.GetPasswordHandler(w, r)
				},
			)
			route.Post(
				"/password/add", func(w http.ResponseWriter, r *http.Request) {
					h.AddPasswordHandler(w, r)
				},
			)
			route.Post(
				"/card/get", func(w http.ResponseWriter, r *http.Request) {
					h.GetCardHandler(w, r)
				},
			)
			route.Post(
				"/card/add", func(w http.ResponseWriter, r *http.Request) {
					h.AddCardHandler(w, r)
				},
			)
		},
	)

	return r
}
