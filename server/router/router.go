package routes

import (
	"context"
	"github.com/egosha7/goph-keeper/internal/compress"
	"github.com/egosha7/goph-keeper/internal/config"
	"github.com/egosha7/goph-keeper/internal/storage"
	"github.com/egosha7/goph-keeper/server/handlers"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"net/http"
)

// SetupRoutes настраивает и возвращает обработчик HTTP-маршрутов.
func SetupRoutes(cfg *config.Config, conn *pgx.Conn, logger *zap.Logger) http.Handler {
	config, err := pgxpool.ParseConfig(cfg.DataBase)
	if err != nil {
		logger.Error("Error parse config", zap.Error(err))
	}
	config.MaxConns = 1000
	// Создание пула подключений
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error("Error connect config", zap.Error(err))
	}

	// Создание хранилища
	// store := storage.NewURLStore(cfg.DataBase, conn, logger, pool)
	repo := storage.NewPostgresURLRepository(conn, logger, pool)

	if cfg.DataBase != "" {
		repo.CreateTable()
	}

	// Создание роутера
	r := chi.NewRouter()

	gzipMiddleware := compress.GzipMiddleware{}

	// Создание группы роутера
	r.Group(
		func(route chi.Router) {
			route.Use(gzipMiddleware.Apply)

			route.Delete(
				"/", func(w http.ResponseWriter, r *http.Request) {

				},
			)

			route.Post(
				"/auth", func(w http.ResponseWriter, r *http.Request) {
					handlers.AuthUser(w, r, repo)
				},
			)

			route.Post(
				"/auth/registration", func(w http.ResponseWriter, r *http.Request) {
					handlers.RegisterUser(w, r, logger, repo)
				},
			)

			route.Post(
				"/pass/namelist", func(w http.ResponseWriter, r *http.Request) {
					handlers.GetSitesList(w, r, logger, repo)
				},
			)

			route.Post(
				"/auth/checkpin", func(w http.ResponseWriter, r *http.Request) {

				},
			)

		},
	)

	return r
}
