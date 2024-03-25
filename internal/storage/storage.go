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

// GetSitesList вставляет информацию о пользователе в таблицу.
func (r *PostgresURLRepository) GetSitesList(login string) ([]string, error) {
	// Использование пула подключений для выполнения запросов
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		r.logger.Error("Error open connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	query := `SELECT name FROM passwords WHERE id_user = (SELECT id FROM users WHERE login = $1)`

	rows, err := r.pool.Query(context.Background(), query, login)
	if err != nil {
		r.logger.Error("Ошибка при выполнении запроса", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var sites []string
	for rows.Next() {
		var site string
		if err := rows.Scan(&site); err != nil {
			r.logger.Error("Ошибка при сканировании результата запроса", zap.Error(err))
			return nil, err
		}
		sites = append(sites, site)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Ошибка в результатах запроса", zap.Error(err))
		return nil, err
	}

	return sites, nil
}

// InsertNewUser вставляет информацию о пользователе в таблицу.
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

// CheckValidUser - валидация пользователя.
func (r *PostgresURLRepository) CheckValidUser(login string) (error, []byte) {
	var pass []byte
	query := "SELECT password FROM users WHERE login = $1"
	err := r.db.QueryRow(context.Background(), query, login).Scan(&pass)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err, nil
		}
		r.logger.Error("Failed to valid user", zap.Error(err))
		return err, nil
	}
	return nil, pass
}

// CheckUniqUser проверяет уникальность логина.
func (r *PostgresURLRepository) CheckUniqUser(login string) (error, string) {
	var existingUser string
	query := "SELECT login FROM users WHERE login = $1"
	err := r.db.QueryRow(context.Background(), query, login).Scan(&existingUser)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err, ""
		}
		r.logger.Error("Failed to check uniq login", zap.Error(err))
		return err, ""
	}
	return nil, existingUser
}

// CreateTable создает необходимые таблицы в базе данных.
func (r *PostgresURLRepository) CreateTable() error {
	_, err := r.db.Exec(
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
