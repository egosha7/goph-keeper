package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/egosha7/goph-keeper/server/internal/domain"
	"github.com/egosha7/goph-keeper/server/internal/service"
	mock_service "github.com/egosha7/goph-keeper/server/internal/service/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_AuthUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, user *domain.User)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          *domain.User
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"login":"Egor","password":"parol","pin":"1234"}`,
			inputUser: &domain.User{
				Login:    "Egor",
				Password: "parol",
				Pin:      "1234",
			},
			mockBehavior: func(s *mock_service.MockServices, user *domain.User) {
				s.EXPECT().RegisterUser(user).Return(nil)
			},
			expectedStatusCode: 200,
		},
	}

	for _, testCase := range testTable {
		t.Run(
			testCase.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				testCase.mockBehavior(auth, testCase.inputUser)
				logger := zap.NewExample()

				// Создаем сервис, используя мок
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/auth", func(w http.ResponseWriter, r *http.Request) {
						handlers.RegisterUser(w, r)
					},
				)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(testCase.inputBody))

				r.ServeHTTP(w, req)

				assert.Equal(t, testCase.expectedStatusCode, w.Code)
			},
		)
	}
}

func TestHandler_AuthenticateUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, user *domain.User)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          *domain.User
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"login":"Egor","password":"parol"}`,
			inputUser: &domain.User{
				Login:    "Egor",
				Password: "parol",
			},
			mockBehavior: func(s *mock_service.MockServices, user *domain.User) {
				s.EXPECT().AuthenticateUser(user).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testTable {
		t.Run(
			testCase.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				testCase.mockBehavior(auth, testCase.inputUser)
				logger := zap.NewExample()

				// Создаем сервис, используя мок
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/auth", func(w http.ResponseWriter, r *http.Request) {
						handlers.AuthUser(w, r)
					},
				)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(testCase.inputBody))

				r.ServeHTTP(w, req)

				assert.Equal(t, testCase.expectedStatusCode, w.Code)
			},
		)
	}
}

func TestHandler_CheckPinCodeHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, login, pin string)

	testCases := []struct {
		name               string
		inputLogin         string
		inputPin           string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:       "Valid Pin",
			inputLogin: "Egor",
			inputPin:   "1234",
			mockBehavior: func(s *mock_service.MockServices, login, pin string) {
				s.EXPECT().CheckPinCode(login, pin).Return(true, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "Invalid Pin",
			inputLogin: "Egor",
			inputPin:   "0000",
			mockBehavior: func(s *mock_service.MockServices, login, pin string) {
				s.EXPECT().CheckPinCode(login, pin).Return(false, nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				tc.mockBehavior(auth, tc.inputLogin, tc.inputPin)
				logger := zap.NewExample()
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/check-pin", func(w http.ResponseWriter, r *http.Request) {
						handlers.CheckPinCodeHandler(w, r)
					},
				)

				// Создаем JSON для запроса
				requestData := domain.CheckPinData{
					Login: tc.inputLogin,
					Pin:   tc.inputPin,
				}
				requestBody, err := json.Marshal(requestData)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/check-pin", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
			},
		)
	}
}
