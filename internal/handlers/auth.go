package handlers

import (
	"encoding/json"
	"github.com/egosha7/goph-keeper/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

// CheckPinCodeHandler обработчик запроса на проверку пин-кода.
func (h *Handler) CheckPinCodeHandler(w http.ResponseWriter, r *http.Request) {
	var requestData domain.CheckPinData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	valid, err := h.Services.CheckPinCode(requestData.Login, requestData.Pin)
	if err != nil {
		h.logger.Error("Ошибка при проверке пин-кода", zap.Error(err))
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

// RegisterUser обрабатывает запрос на регистрацию нового пользователя.
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	if err := h.Services.RegisterUser(&user); err != nil {
		h.logger.Error("Ошибка при регистрации пользователя", zap.Error(err))
		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AuthUser обрабатывает запрос аутентификации пользователя.
func (h *Handler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var user *domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	h.logger.Info(user.Password + user.Login + user.Pin)

	if err := h.Services.AuthenticateUser(user); err != nil {
		http.Error(w, "Неверная пара логин/пароль", http.StatusUnauthorized)
		h.logger.Error("Failed to check user validity", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
