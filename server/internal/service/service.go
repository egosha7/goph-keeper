package service

import (
	"github.com/egosha7/goph-keeper/server/internal/domain"
	"github.com/egosha7/goph-keeper/server/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Services представляет сервис для работы с пользователями.
//
//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Services interface {
	CheckPinCode(login, pin string) (bool, error)
	AddPassword(login, passName, password string) error
	GetPassword(login, passName string) (string, error)
	AddCard(login, cardName, numberCard, expiryDateCard, cvvCard string) error
	GetCard(login, cardName string) (string, string, string, error)
	GetPasswordNameList(login string) ([]string, error)
	GetCardNameList(login string) ([]string, error)
	RegisterUser(user *domain.User) error
	AuthenticateUser(user *domain.User) error
}

type Service struct {
	Services
}

// UserServiceImpl представляет реализацию UserService.
type UserServiceImpl struct {
	Repository repository.UserRepository
	Logger     *zap.Logger
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(repository *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		Services: &UserServiceImpl{
			Repository: repository,
			Logger:     logger,
		},
	}
}

// CheckPinCode проверяет пин-код для указанного пользователя.
func (s *UserServiceImpl) CheckPinCode(login, pin string) (bool, error) {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.CheckPinCode(login, pin)
}

// AddPassword добавляет новый пароль.
func (s *UserServiceImpl) AddPassword(login, passName, password string) error {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.InsertNewPassword(login, passName, password)
}

// GetPassword возвращает пароль по его имени.
func (s *UserServiceImpl) GetPassword(login, passName string) (string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.GetPassword(login, passName)
}

// AddCard добавляет новую карту.
func (s *UserServiceImpl) AddCard(login, cardName, numberCard, expiryDateCard, cvvCard string) error {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.InsertNewCard(login, cardName, numberCard, expiryDateCard, cvvCard)
}

// GetCard возвращает информацию о карте по ее имени.
func (s *UserServiceImpl) GetCard(login, cardName string) (string, string, string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.GetCard(login, cardName)
}

// GetPasswordNameList возвращает список названий паролей для указанного пользователя.
func (s *UserServiceImpl) GetPasswordNameList(login string) ([]string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.GetPasswordNameList(login)
}

// GetCardNameList возвращает список названий карт для указанного пользователя.
func (s *UserServiceImpl) GetCardNameList(login string) ([]string, error) {
	// Здесь может быть ваша бизнес-логика
	return s.Repository.GetCardNameList(login)
}

// RegisterUser регистрирует нового пользователя.
func (s *UserServiceImpl) RegisterUser(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.Repository.Create(user)
}

// AuthenticateUser аутентифицирует пользователя.
func (s *UserServiceImpl) AuthenticateUser(user *domain.User) error {
	storedUser, err := s.Repository.CheckValidUser(user.Login)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(storedUser), []byte(user.Password))
}
