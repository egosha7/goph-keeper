package service

import (
	"github.com/egosha7/goph-keeper/server/internal/domain"
	"github.com/egosha7/goph-keeper/server/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService представляет сервис для работы с пользователями.
type UserService struct {
	repository *repository.PostgreSQLRepository
	logger     *zap.Logger
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(repository *repository.PostgreSQLRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
	}
}

// CheckPinCode проверяет пин-код для указанного пользователя.
func (s *UserService) CheckPinCode(login, pin string) (bool, error) {
	// Здесь может быть ваша бизнес-логика
	return s.repository.CheckPinCode(login, pin)
}

// AddPassword добавляет новый пароль.
func (s *UserService) AddPassword(login, passName, password string) error {
	// Здесь может быть ваша бизнес-логика
	return s.repository.InsertNewPassword(login, passName, password)
}

// GetPassword возвращает пароль по его имени.
func (s *UserService) GetPassword(login, passName string) (string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.repository.GetPassword(login, passName)
}

// AddCard добавляет новую карту.
func (s *UserService) AddCard(login, cardName, numberCard, expiryDateCard, cvvCard string) error {
	// Здесь может быть ваша бизнес-логика
	return s.repository.InsertNewCard(login, cardName, numberCard, expiryDateCard, cvvCard)
}

// GetCard возвращает информацию о карте по ее имени.
func (s *UserService) GetCard(login, cardName string) (string, string, string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.repository.GetCard(login, cardName)
}

// GetPasswordNameList возвращает список названий паролей для указанного пользователя.
func (s *UserService) GetPasswordNameList(login string) ([]string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.repository.GetPasswordNameList(login)
}

// GetCardNameList возвращает список названий карт для указанного пользователя.
func (s *UserService) GetCardNameList(login string) ([]string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.repository.GetCardNameList(login)
}

// RegisterUser регистрирует нового пользователя.
func (s *UserService) RegisterUser(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repository.Create(user)
}

// AuthenticateUser аутентифицирует пользователя.
func (s *UserService) AuthenticateUser(user *domain.User) error {
	storedUser, err := s.repository.GetByUsername(user.Login)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
}
