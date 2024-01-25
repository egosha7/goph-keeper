package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
	"sync"
)

// URLStore представляет хранилище сокращенных URL.
type URLStore struct {
	urls     []URL
	mu       sync.RWMutex
	DBstring string
	db       *pgx.Conn
	logger   *zap.Logger
	pool     *pgxpool.Pool
}

// URL представляет структуру с данными о сокращенном URL.
type URL struct {
	ID     string
	URL    string
	UserID string
}

// NewURLStore создает новый экземпляр URLStore.
func NewURLStore(DBstring string, db *pgx.Conn, logger *zap.Logger, pool *pgxpool.Pool) *URLStore {
	return &URLStore{
		urls:     make([]URL, 0),
		DBstring: DBstring,
		db:       db,
		logger:   logger,
		pool:     pool,
	}
}

// URLRepository представляет интерфейс для работы с базой данных.
type URLRepository interface {
	AddURL(id string, url string) (string, bool)
	GetIDByURL(url string) (string, bool)
	GetURLByID(id string) (string, bool)
	CreateTable()
	PrintAllURLs()
}

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
