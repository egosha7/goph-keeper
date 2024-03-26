package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/egosha7/goph-keeper/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

// User структура для представления пользователя
type UserInfo struct {
	Login string `json:"Login"`
}

// CheckPinData структура для передачи данных при проверке пин-кода.
type CheckPinData struct {
	Login string `json:"login"`
	Pin   string `json:"pin"`
}

// CheckPinResponse структура для ответа о результате проверки пин-кода.
type CheckPinResponse struct {
	Valid bool `json:"valid"`
}

type PassData struct {
	Login    string `json:"login"`
	PassName string `json:"passName"`
}

type CardData struct {
	Login    string `json:"login"`
	CardName string `json:"cardName"`
}

// PassData содержит информацию о новом пароле.
type PasswordData struct {
	Login    string `json:"login"`    // Логин пользователя
	PassName string `json:"passName"` // Название нового пароля
	Password string `json:"password"` // Новый пароль
}

// NewCardData содержит информацию о новой карте.
type NewCardData struct {
	Login          string `json:"login"`          // Логин пользователя
	CardName       string `json:"cardName"`       // Название новой карты
	NumberCard     string `json:"numberCard"`     // Номер новой карты
	ExpiryDateCard string `json:"expiryDateCard"` // Срок новой карты
	CvvCard        string `json:"CvvCard"`        // Секретный код новой карты
}

func NewPassword(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Парсинг JSON-данных из запроса
	var passwordData PasswordData
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Вставка нового пользователя в таблицу users и получение user_id
	err = store.InsertNewPassword(passwordData.Login, passwordData.PassName, passwordData.Password)
	if err != nil {
		logger.Error("Ошибка вставки нового пароля", zap.Error(err))
		http.Error(w, "Ошибка вставки нового пароля", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	logger.Info("Пароль успешно сохранен", zap.String("login", passwordData.PassName))
	fmt.Fprintf(w, "Пароль %s успешно сохранен", passwordData.PassName)
}

func NewCard(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Парсинг JSON-данных из запроса
	var newCardData NewCardData
	err := json.NewDecoder(r.Body).Decode(&newCardData)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Вставка нового пользователя в таблицу users и получение user_id
	err = store.InsertNewCard(newCardData.Login, newCardData.CardName, newCardData.NumberCard, newCardData.ExpiryDateCard, newCardData.CvvCard)
	if err != nil {
		logger.Error("Ошибка вставки новой карты", zap.Error(err))
		http.Error(w, "Ошибка вставки новой карты", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	logger.Info("Карта успешно сохранена", zap.String("login", newCardData.CardName))
	fmt.Fprintf(w, "Карта %s успешно сохранена", newCardData.CardName)
}

func GetPasswordNameList(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Получаем логин пользователя из запроса
	var userInfo UserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Получаем список сайтов из базы данных для данного пользователя
	passwords, err := store.GetPasswordNameList(userInfo.Login)
	if err != nil {
		logger.Error("Ошибка при получении списка сайтов", zap.Error(err))
		http.Error(w, "Ошибка при получении списка сайтов", http.StatusInternalServerError)
		return
	}

	// Отправляем список сайтов клиенту в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(passwords); err != nil {
		logger.Error("Ошибка при кодировании данных в JSON", zap.Error(err))
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}
}

func GetCardList(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Получаем логин пользователя из запроса
	var userInfo UserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Получаем список сайтов из базы данных для данного пользователя
	cards, err := store.GetCardList(userInfo.Login)
	if err != nil {
		logger.Error("Ошибка при получении списка сайтов", zap.Error(err))
		http.Error(w, "Ошибка при получении списка сайтов", http.StatusInternalServerError)
		return
	}

	// Отправляем список сайтов клиенту в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		logger.Error("Ошибка при кодировании данных в JSON", zap.Error(err))
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}
}

// CheckPinCodeHandler обработчик запроса на проверку пин-кода.
func CheckPinCodeHandler(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	var requestData CheckPinData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	valid, err := store.CheckPinCode(requestData.Login, requestData.Pin)
	if err != nil {
		logger.Error("Ошибка при проверке пин-кода", zap.Error(err))
		http.Error(w, "Ошибка при проверке пин-кода", http.StatusInternalServerError)
		return
	}

	if valid {
		// Если пин-коды совпадают, отправляем статус OK
		w.WriteHeader(http.StatusOK)
		return
	}

	// Если пин-коды не совпадают, отправляем статус Unauthorized
	http.Error(w, "Неверный пин-код", http.StatusUnauthorized)
}

func GetPasswordHandler(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Парсим JSON-данные из запроса
	var passData PassData
	err := json.NewDecoder(r.Body).Decode(&passData)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Получаем пароль из базы данных по логину и названию пароля
	password, err := store.GetPassword(passData.Login, passData.PassName)
	if err != nil {
		logger.Error("Ошибка при получении пароля", zap.Error(err))
		http.Error(w, fmt.Sprintf("ошибка при получении пароля: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем пароль клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, password)
}

// CardInfo содержит информацию о банковской карте.
type CardInfo struct {
	Number     string `json:"number"`     // Номер карты
	ExpiryDate string `json:"expiryDate"` // Срок действия карты
	CVV        string `json:"cvv"`        // CVV карты
}

func GetCardHandler(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Парсим JSON-данные из запроса
	var cardData CardData
	err := json.NewDecoder(r.Body).Decode(&cardData)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Получаем данные карты из базы данных по логину и названию карты
	cardNumber, cardExpiryDate, cardCVV, err := store.GetCard(cardData.Login, cardData.CardName)
	if err != nil {
		logger.Error("Ошибка при получении данных о карте", zap.Error(err))
		http.Error(w, fmt.Sprintf("ошибка при получении данных о карте: %v", err), http.StatusInternalServerError)
		return
	}

	// Формируем структуру с данными о карте
	cardInfo := CardInfo{
		Number:     cardNumber,
		ExpiryDate: cardExpiryDate,
		CVV:        cardCVV,
	}

	// Кодируем данные о карте в JSON
	jsonData, err := json.Marshal(cardInfo)
	if err != nil {
		logger.Error("Ошибка при кодировании данных о карте в JSON", zap.Error(err))
		http.Error(w, "Ошибка при кодировании данных о карте в JSON", http.StatusInternalServerError)
		return
	}

	// Отправляем данные о карте клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		logger.Error("Ошибка при отправке данных о карте клиенту", zap.Error(err))
	}
}
