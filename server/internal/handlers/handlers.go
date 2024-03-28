package handlers

import (
	"github.com/egosha7/goph-keeper/server/internal/service"
	"go.uber.org/zap"
)

// Handler представляет обработчик HTTP-запросов.
type Handler struct {
	Services *service.Service
	logger   *zap.Logger
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(services *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		Services: services,
		logger:   logger,
	}
}
