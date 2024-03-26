package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/egosha7/goph-keeper/internal/storage"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// User структура для представления пользователя
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Pin      string `json:"pin"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Парсинг JSON-данных из запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей в запросе
	if user.Login == "" || user.Password == "" || user.Pin == "" {
		logger.Error("Неверный формат запроса: отсутствуют обязательные поля")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка, что логин уникальный
	var existingUser string
	err, existingUser = store.CheckUniqUser(user.Login)
	if err != nil && err != pgx.ErrNoRows {
		logger.Error("Ошибка при выполнении запроса к базе данных", zap.Error(err))
		http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
		return
	}
	if existingUser != "" {
		logger.Info("Пользователь с таким логином уже существует", zap.String("login", user.Login))
		http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
		return
	}

	// Хэширование пароля перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Ошибка хэширования пароля", zap.Error(err))
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
		return
	}

	// Вставка нового пользователя в таблицу users
	err = store.InsertNewUser(user.Login, hashedPassword, user.Pin)
	if err != nil {
		logger.Error("Ошибка вставки пользователя", zap.Error(err))
		http.Error(w, "Ошибка вставки пользователя", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	logger.Info("Пользователь успешно зарегистрирован", zap.String("login", user.Login))
	fmt.Fprintf(w, "Пользователь %s успешно аутентифицирован", user.Login)
}

func AuthUser(w http.ResponseWriter, r *http.Request, store *storage.PostgresURLRepository) {
	// Парсинг JSON-данных из запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей в запросе
	if user.Login == "" || user.Password == "" {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Получение хэшированного пароля из базы данных по логину
	var hashedPassword []byte
	err, hashedPassword = store.CheckValidUser(user.Login)
	if err == pgx.ErrNoRows {
		http.Error(w, "Неверная пара логин/пароль", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Неверная пара логин/пароль", http.StatusUnauthorized)
		return
	}

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Пользователь %s успешно аутентифицирован", user.Login)
}
