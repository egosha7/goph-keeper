package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/egosha7/goph-keeper/server/internal/domain"
	"github.com/egosha7/goph-keeper/server/internal/service"
	mock_service "github.com/egosha7/goph-keeper/server/internal/service/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_AddPasswordHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, requestData domain.PasswordData)

	testCases := []struct {
		name                 string
		requestData          domain.PasswordData
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid Password",
			requestData: domain.PasswordData{
				Login:    "Egor",
				PassName: "Email",
				Password: "password123",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain.PasswordData) {
				s.EXPECT().AddPassword(requestData.Login, requestData.PassName, requestData.Password).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Password",
			requestData: domain.PasswordData{
				Login:    "NonExistentUser",
				PassName: "Email",
				Password: "password123",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain.PasswordData) {
				s.EXPECT().AddPassword(
					requestData.Login, requestData.PassName, requestData.Password,
				).Return(errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `Ошибка при добавлении нового пароля`,
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
					"/add-password", func(w http.ResponseWriter, r *http.Request) {
						handlers.AddPasswordHandler(w, r)
					},
				)

				requestBody, err := json.Marshal(tc.requestData)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/add-password", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
				if tc.expectedResponseBody != "" {
					assert.Equal(t, tc.expectedResponseBody, strings.TrimSpace(w.Body.String()))
				}
			},
		)
	}
}

func TestHandler_GetPasswordHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, requestData domain.PassData)

	testCases := []struct {
		name                 string
		requestData          domain.PassData
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid Password",
			requestData: domain.PassData{
				Login:    "Egor",
				PassName: "Email",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain.PassData) {
				s.EXPECT().GetPassword(requestData.Login, requestData.PassName).Return("password123", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "password123",
		},
		{
			name: "Invalid Password",
			requestData: domain.PassData{
				Login:    "NonExistentUser",
				PassName: "Email",
			},
			mockBehavior: func(s *mock_service.MockServices, requestData domain.PassData) {
				s.EXPECT().GetPassword(requestData.Login, requestData.PassName).Return("", errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `Ошибка при получении пароля`,
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
					"/get-password", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetPasswordHandler(w, r)
					},
				)

				requestBody, err := json.Marshal(tc.requestData)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/get-password", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
				if tc.expectedResponseBody != "" {
					assert.Equal(t, tc.expectedResponseBody, strings.TrimSpace(w.Body.String()))
				}
			},
		)
	}
}

func TestHandler_GetPasswordNameList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockServices, userInfo domain.UserInfo)

	testCases := []struct {
		name                 string
		userInfo             domain.UserInfo
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid User",
			userInfo: domain.UserInfo{
				Login: "Egor",
			},
			mockBehavior: func(s *mock_service.MockServices, userInfo domain.UserInfo) {
				s.EXPECT().GetPasswordNameList(userInfo.Login).Return([]string{"Email", "SocialMedia"}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `["Email","SocialMedia"]`,
		},
		{
			name: "Invalid User",
			userInfo: domain.UserInfo{
				Login: "NonExistentUser",
			},
			mockBehavior: func(s *mock_service.MockServices, userInfo domain.UserInfo) {
				s.EXPECT().GetPasswordNameList(userInfo.Login).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `Ошибка при получении списка названий паролей`,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				auth := mock_service.NewMockServices(ctrl)
				tc.mockBehavior(auth, tc.userInfo)
				logger := zap.NewExample()
				userService := &service.Service{Services: auth}

				handlers := NewHandler(userService, logger)

				r := chi.NewRouter()
				r.Post(
					"/get-password-name-list", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetPasswordNameList(w, r)
					},
				)

				requestBody, err := json.Marshal(tc.userInfo)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/get-password-name-list", bytes.NewBuffer(requestBody))
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatusCode, w.Code)
				if tc.expectedResponseBody != "" {
					assert.Equal(t, tc.expectedResponseBody, strings.TrimSpace(w.Body.String()))
				}
			},
		)
	}
}
