package db

import (
	"context"
	"fmt"
	"github.com/egosha7/goph-keeper/internal/config"
	"github.com/jackc/pgx/v4"
	"net/http"
)

// ConnectToDB устанавливает соединение с базой данных на основе конфигурации.
// Возвращает соединение (pgx.Conn) и ошибку, если возникает ошибка при подключении.
func ConnectToDB(cfg *config.Config) (*pgx.Conn, error) {
	if cfg.DataBase == "" {
		// Возвращаем пустое соединение, если строка подключения пуста
		conn := &pgx.Conn{}
		return conn, nil
	}

	connConfig, err := pgx.ParseConfig(cfg.DataBase)
	if err != nil {
		return nil, err
	}

	fmt.Printf(cfg.DataBase)

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// PingDB выполняет пинг базы данных и отправляет статус в HTTP-ответ.
func PingDB(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	err := conn.Ping(context.Background())
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
