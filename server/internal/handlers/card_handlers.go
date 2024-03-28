package handlers

import (
	"encoding/json"
	"github.com/egosha7/goph-keeper/server/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

// AddCardHandler обрабатывает запрос на добавление новой карты.
func (h *Handler) AddCardHandler(w http.ResponseWriter, r *http.Request) {
	var requestData domain.NewCardData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	err := h.Services.AddCard(
		requestData.Login, requestData.CardName, requestData.NumberCard, requestData.ExpiryDateCard,
		requestData.CvvCard,
	)
	if err != nil {
		h.logger.Error("Ошибка при добавлении новой карты", zap.Error(err))
		http.Error(w, "Ошибка при добавлении новой карты", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetCardHandler обрабатывает запрос на получение информации о карте.
func (h *Handler) GetCardHandler(w http.ResponseWriter, r *http.Request) {
	var requestData domain.CardData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	cardNumber, cardExpiryDate, cardCVV, err := h.Services.GetCard(requestData.Login, requestData.CardName)
	if err != nil {
		h.logger.Error("Ошибка при получении информации о карте", zap.Error(err))
		http.Error(w, "Ошибка при получении информации о карте", http.StatusInternalServerError)
		return
	}

	type CardInfo struct {
		Number     string `json:"number"`
		ExpiryDate string `json:"expiryDate"`
		CVV        string `json:"cvv"`
	}

	cardInfo := CardInfo{
		Number:     cardNumber,
		ExpiryDate: cardExpiryDate,
		CVV:        cardCVV,
	}

	response, err := json.Marshal(cardInfo)
	if err != nil {
		h.logger.Error("Ошибка при кодировании данных в JSON", zap.Error(err))
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetCardList обрабатывает запрос на получение списка названий карт для указанного пользователя.
func (h *Handler) GetCardList(w http.ResponseWriter, r *http.Request) {
	var userInfo domain.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		h.logger.Error("Ошибка при разборе JSON", zap.Error(err))
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	cards, err := h.Services.GetCardNameList(userInfo.Login)
	if err != nil {
		h.logger.Error("Ошибка при получении списка названий карт", zap.Error(err))
		http.Error(w, "Ошибка при получении списка названий карт", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		h.logger.Error("Ошибка при кодировании данных в JSON", zap.Error(err))
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}
}
