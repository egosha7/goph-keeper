package handlers

import (
	"encoding/json"
	"github.com/egosha7/goph-keeper/server/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

// AddPasswordHandler обрабатывает запрос на добавление нового пароля.
func (h *Handler) AddPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var requestData domain.PasswordData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	err := h.Services.AddPassword(requestData.Login, requestData.PassName, requestData.Password)
	if err != nil {
		h.logger.Error("Ошибка при добавлении нового пароля", zap.Error(err))
		http.Error(w, "Ошибка при добавлении нового пароля", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetPasswordHandler обрабатывает запрос на получение пароля.
func (h *Handler) GetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var requestData domain.PassData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	password, err := h.Services.GetPassword(requestData.Login, requestData.PassName)
	if err != nil {
		h.logger.Error("Ошибка при получении пароля", zap.Error(err))
		http.Error(w, "Ошибка при получении пароля", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(password))
}

// GetPasswordNameList обрабатывает запрос на получение списка названий паролей для указанного пользователя.
func (h *Handler) GetPasswordNameList(w http.ResponseWriter, r *http.Request) {
	var userInfo domain.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	passwords, err := h.Services.GetPasswordNameList(userInfo.Login)
	if err != nil {
		h.logger.Error("Ошибка при получении списка названий паролей", zap.Error(err))
		http.Error(w, "Ошибка при получении списка названий паролей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(passwords); err != nil {
		h.logger.Error("Ошибка при кодировании данных в JSON", zap.Error(err))
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}
}
