package handlers

import (
	"encoding/json"
	"github.com/egosha7/goph-keeper/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

// User структура для представления пользователя
type UserInfo struct {
	Login string `json:"Login"`
}

func GetSitesList(w http.ResponseWriter, r *http.Request, logger *zap.Logger, store *storage.PostgresURLRepository) {
	// Получаем логин пользователя из запроса
	var userInfo UserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Получаем список сайтов из базы данных для данного пользователя
	passwords, err := store.GetSitesList(userInfo.Login)
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
