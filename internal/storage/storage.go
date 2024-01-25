package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

// PostgresURLRepository реализует интерфейс URLRepository для работы с PostgreSQL.
type PostgresURLRepository struct {
	db     *pgx.Conn
	logger *zap.Logger
	pool   *pgxpool.Pool
}

// NewPostgresURLRepository создает новый экземпляр PostgresURLRepository.
func NewPostgresURLRepository(db *pgx.Conn, logger *zap.Logger, pool *pgxpool.Pool) *PostgresURLRepository {
	return &PostgresURLRepository{
		db:     db,
		logger: logger,
		pool:   pool,
	}
}

// CreateTable создает необходимые таблицы в базе данных.
func (r *PostgresURLRepository) InsertNewUser(login string, password []byte) error {
	// Использование пула подключений для выполнения запросов
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		r.logger.Error("Error open connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	// Добавляем данные в таблицу user_urls
	userQuery := "INSERT INTO users (login, password) VALUES ($1, $2)"
	_, userErr := conn.Exec(context.Background(), userQuery, login, password)
	if userErr != nil {
		r.logger.Error("Failed to add user URL", zap.Error(userErr))
		conn.Release()
		return userErr
	}
	conn.Release()
	return nil
}

// CreateTable создает необходимые таблицы в базе данных.
func (r *PostgresURLRepository) CheckValidUser(login string) (error, string) {
	var pass string
	query := "SELECT password FROM users WHERE login = $1"
	err := r.db.QueryRow(context.Background(), query, login).Scan(&pass)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err, ""
		}
		r.logger.Error("Failed to get ID by URL", zap.Error(err))
		return err, ""
	}
	return nil, pass
}

// CreateTable создает необходимые таблицы в базе данных.
func (r *PostgresURLRepository) CreateTable() error {
	_, err := r.db.Exec(
		context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			ID TEXT PRIMARY KEY,
			URL TEXT,
			UNIQUE (URL)
		)
	`,
	)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		context.Background(), `
		CREATE TABLE IF NOT EXISTS user_urls (
			ID SERIAL PRIMARY KEY,
			IDshortURL TEXT,
			userID TEXT,
			delFLAG BOOL DEFAULT false
		)
	`,
	)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		context.Background(), `
		ALTER TABLE user_urls
		ADD CONSTRAINT fk_name_IDshortURL
		FOREIGN KEY (IDshortURL) REFERENCES urls (ID);

	`,
	)
	if err != nil {
		return err
	}

	return nil
}
