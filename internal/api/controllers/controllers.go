package controllers

import (
	"github.com/SimpleOG/WebSocketChat/internal/api/controllers/RoomControllers"
	"github.com/SimpleOG/WebSocketChat/internal/api/controllers/UserControllers"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/service"
)

type Controllers struct {
	UserControllers.UserControllers
	RoomControllers.RoomControllers
}

func NewControllers(logger logger.Logger, service *service.Service) *Controllers {
	return &Controllers{
		UserControllers: UserControllers.NewUserControllers(logger, service),
		RoomControllers: RoomControllers.NewRoomControllers(logger, service),
	}
}
