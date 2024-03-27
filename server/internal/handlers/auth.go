package handlers

import (
	"encoding/json"
	"github.com/egosha7/goph-keeper/server/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

// RegisterUser обрабатывает запрос на регистрацию нового пользователя.
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	if err := h.userService.RegisterUser(&user); err != nil {
		h.logger.Error("Ошибка при регистрации пользователя", zap.Error(err))
		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AuthUser обрабатывает запрос аутентификации пользователя.
func (h *Handler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	if err := h.userService.AuthenticateUser(&user); err != nil {
		http.Error(w, "Неверная пара логин/пароль", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
