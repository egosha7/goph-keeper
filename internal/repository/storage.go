package repository

import (
	"context"
	"fmt"
	"github.com/egosha7/goph-keeper/internal/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// UserRepository представляет интерфейс для работы с данными пользователей.
type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	CheckUniqUser(login string) (bool, error)
	CheckPinCode(login, pin string) (bool, error)
	CheckValidUser(login string) (string, error)
	InsertNewCard(login, cardName, numberCard, expiryDateCard, cvvCard string) error
	GetCard(login, cardName string) (string, string, string, error)
	GetCardNameList(login string) ([]string, error)
	InsertNewPassword(login, passName, password string) error
	GetPassword(login, passName string) (string, error)
	GetPasswordNameList(login string) ([]string, error)
}

type Repository struct {
	UserRepository
}

// PostgreSQLRepository представляет репозиторий для работы с PostgreSQL.
type PostgreSQLRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgreSQLRepository создает новый экземпляр PostgreSQLRepository.
func NewPostgreSQLRepository(pool *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{
		UserRepository: &PostgreSQLRepository{
			pool:   pool,
			logger: logger,
		},
	}
}

// CheckValidUser проверяет валидность пользователя.
func (r *PostgreSQLRepository) CheckValidUser(login string) (string, error) {
	var password string
	query := "SELECT password FROM users WHERE login = $1"
	err := r.pool.QueryRow(context.Background(), query, login).Scan(&password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		r.logger.Error("Failed to check user validity", zap.Error(err))
		return "", err
	}
	return password, nil
}

// CheckUniqUser проверяет уникальность логина пользователя.
func (r *PostgreSQLRepository) CheckUniqUser(login string) (bool, error) {
	var existingUser string
	query := "SELECT login FROM users WHERE login = $1"
	err := r.pool.QueryRow(context.Background(), query, login).Scan(&existingUser)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		r.logger.Error("Failed to check unique user", zap.Error(err))
		return false, err
	}
	return true, nil
}

// InsertNewCard вставляет новую карту в базу данных.
func (r *PostgreSQLRepository) InsertNewCard(login, cardName, numberCard, expiryDateCard, cvvCard string) error {
	query := "INSERT INTO cards (number, expirydate, cvv, id_user, name) VALUES ($1, $2, $3, (SELECT id FROM users WHERE login = $4), $5)"
	_, err := r.pool.Exec(context.Background(), query, numberCard, expiryDateCard, cvvCard, login, cardName)
	if err != nil {
		r.logger.Error("Failed to insert new card", zap.Error(err))
		return err
	}
	return nil
}

// GetCard получает данные о карте из базы данных.
func (r *PostgreSQLRepository) GetCard(login, cardName string) (string, string, string, error) {
	var cardNumber, cardExpiryDate, cardCVV string
	query := `SELECT number, expirydate, cvv FROM cards WHERE id_user = (SELECT id FROM users WHERE login = $1) AND name = $2`
	err := r.pool.QueryRow(context.Background(), query, login, cardName).Scan(&cardNumber, &cardExpiryDate, &cardCVV)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", fmt.Errorf("card not found")
		}
		r.logger.Error("Failed to get card", zap.Error(err))
		return "", "", "", err
	}
	return cardNumber, cardExpiryDate, cardCVV, nil
}

// GetCardList получает список имен карт пользователя из базы данных.
func (r *PostgreSQLRepository) GetCardNameList(login string) ([]string, error) {
	query := `SELECT name FROM cards WHERE id_user = (SELECT id FROM users WHERE login = $1)`
	rows, err := r.pool.Query(context.Background(), query, login)
	if err != nil {
		r.logger.Error("Failed to get card list", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var cardNames []string
	for rows.Next() {
		var cardName string
		if err := rows.Scan(&cardName); err != nil {
			r.logger.Error("Failed to scan card list row", zap.Error(err))
			return nil, err
		}
		cardNames = append(cardNames, cardName)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error in card list rows", zap.Error(err))
		return nil, err
	}
	return cardNames, nil
}

// InsertNewPassword вставляет новый пароль в базу данных.
func (r *PostgreSQLRepository) InsertNewPassword(login, passName, password string) error {
	query := "INSERT INTO passwords (id_user, name, password) VALUES ((SELECT id FROM users WHERE login = $1), $2, $3)"
	_, err := r.pool.Exec(context.Background(), query, login, passName, password)
	if err != nil {
		r.logger.Error("Failed to insert new password", zap.Error(err))
		return err
	}
	return nil
}

// GetPassword получает пароль из базы данных.
func (r *PostgreSQLRepository) GetPassword(login, passName string) (string, error) {
	var password string
	query := `SELECT password FROM passwords WHERE id_user = (SELECT id FROM users WHERE login = $1) AND name = $2`
	err := r.pool.QueryRow(context.Background(), query, login, passName).Scan(&password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("password not found")
		}
		r.logger.Error("Failed to get password", zap.Error(err))
		return "", err
	}
	return password, nil
}

// GetPasswordNameList получает список имен паролей пользователя из базы данных.
func (r *PostgreSQLRepository) GetPasswordNameList(login string) ([]string, error) {
	query := `SELECT name FROM passwords WHERE id_user = (SELECT id FROM users WHERE login = $1)`
	rows, err := r.pool.Query(context.Background(), query, login)
	if err != nil {
		r.logger.Error("Failed to get password name list", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var passNames []string
	for rows.Next() {
		var passName string
		if err := rows.Scan(&passName); err != nil {
			r.logger.Error("Failed to scan password name list row", zap.Error(err))
			return nil, err
		}
		passNames = append(passNames, passName)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error in password name list rows", zap.Error(err))
		return nil, err
	}
	return passNames, nil
}

// CheckPinCode проверяет пин-код пользователя.
func (r *PostgreSQLRepository) CheckPinCode(login, pin string) (bool, error) {
	var valid bool
	query := `SELECT EXISTS (SELECT * FROM users WHERE login = $1 AND pin = $2)`
	err := r.pool.QueryRow(context.Background(), query, login, pin).Scan(&valid)
	if err != nil {
		r.logger.Error("Failed to check pin code", zap.Error(err))
		return false, err
	}
	return valid, nil
}

// Create создает нового пользователя в базе данных.
func (r *PostgreSQLRepository) Create(user *domain.User) error {
	_, err := r.pool.Exec(
		context.Background(), "INSERT INTO users (login, password, pin) VALUES ($1, $2, $3)", user.Login, user.Password,
		user.Pin,
	)
	return err
}

// GetByUsername возвращает пользователя из базы данных по его логину.
func (r *PostgreSQLRepository) GetByUsername(username string) (*domain.User, error) {
	row := r.pool.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", username)
	user := &domain.User{}
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
