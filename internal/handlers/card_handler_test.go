package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	domain2 "github.com/egosha7/goph-keeper/internal/domain"
	"github.com/egosha7/goph-keeper/internal/service"
	mock_service "github.com/egosha7/goph-keeper/internal/service/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_AddCardHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, requestData domain2.NewCardData)

	testCases := []struct {
		name               string
		requestData        domain2.NewCardData
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "Valid Card",
			requestData: domain2.NewCardData{
				Login:          "Egor",
				CardName:       "Visa",
				NumberCard:     "1234567890123456",
				ExpiryDateCard: "12/24",
				CvvCard:        "123",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain2.NewCardData) {
				s.EXPECT().AddCard(
					requestData.Login, requestData.CardName, requestData.NumberCard, requestData.ExpiryDateCard,
					requestData.CvvCard,
				).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Card",
			requestData: domain2.NewCardData{
				Login:          "Egor",
				CardName:       "Mastercard",
				NumberCard:     "123456789012345",
				ExpiryDateCard: "12/25",
				CvvCard:        "1234",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain2.NewCardData) {
				s.EXPECT().AddCard(
					requestData.Login, requestData.CardName, requestData.NumberCard, requestData.ExpiryDateCard,
					requestData.CvvCard,
				).Return(errors.New("error adding card"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				tc.mockBehavior(auth, tc.requestData)
				logger := zap.NewExample()
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/add-card", func(w http.ResponseWriter, r *http.Request) {
						handlers.AddCardHandler(w, r)
					},
				)

				requestBody, err := json.Marshal(tc.requestData)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/add-card", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
			},
		)
	}
}

func TestHandler_GetCardHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, requestData domain2.CardData)

	testCases := []struct {
		name                 string
		requestData          domain2.CardData
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid Card",
			requestData: domain2.CardData{
				Login:    "Egor",
				CardName: "Visa",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain2.CardData) {
				s.EXPECT().GetCard(requestData.Login, requestData.CardName).Return(
					"1234567890123456", "12/24", "123", nil,
				)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"number":"1234567890123456","expiryDate":"12/24","cvv":"123"}`,
		},
		{
			name: "Invalid Card",
			requestData: domain2.CardData{
				Login:    "Egor",
				CardName: "Mastercard",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain2.CardData) {
				s.EXPECT().GetCard(requestData.Login, requestData.CardName).Return(
					"", "", "", errors.New("card not found"),
				)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `Ошибка при получении информации о карте`,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				tc.mockBehavior(auth, tc.requestData)
				logger := zap.NewExample()
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/get-card", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetCardHandler(w, r)
					},
				)

				requestBody, err := json.Marshal(tc.requestData)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/get-card", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
				assert.Equal(t, tc.expectedResponseBody, strings.TrimSpace(w.Body.String()))
			},
		)
	}
}

func TestHandler_GetCardList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, userInfo domain2.UserInfo)

	testCases := []struct {
		name                 string
		requestData          domain2.UserInfo
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid User",
			requestData: domain2.UserInfo{
				Login: "Egor",
			},
			mockBehavior: func(s *mock_service.MockServices, userInfo domain2.UserInfo) {
				s.EXPECT().GetCardNameList(userInfo.Login).Return([]string{"Visa", "Mastercard"}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `["Visa","Mastercard"]`,
		},
		{
			name: "Invalid User",
			requestData: domain2.UserInfo{
				Login: "NonExistentUser",
			},
			mockBehavior: func(s *mock_service.MockServices, userInfo domain2.UserInfo) {
				s.EXPECT().GetCardNameList(userInfo.Login).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `Ошибка при получении списка названий карт`,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				tc.mockBehavior(auth, tc.requestData)
				logger := zap.NewExample()
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/get-card-list", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetCardList(w, r)
					},
				)

				requestBody, err := json.Marshal(tc.requestData)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/get-card-list", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
				assert.Equal(
					t, tc.expectedResponseBody, strings.TrimSpace(w.Body.String()),
				) // Убираем пробельные символы
			},
		)
	}
}
