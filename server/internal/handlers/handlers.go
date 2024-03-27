package handlers

import (
	"github.com/egosha7/goph-keeper/server/internal/service"
	"go.uber.org/zap"
)

// Handler представляет обработчик HTTP-запросов.
type Handler struct {
	userService *service.UserService
	logger      *zap.Logger
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(userService *service.UserService, logger *zap.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}
