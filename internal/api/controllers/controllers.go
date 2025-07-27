package controllers

import (
	"github.com/SimpleOG/WebSocketChat/internal/api/controllers/RoomControllers"
	"github.com/SimpleOG/WebSocketChat/internal/api/controllers/UserControllers"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	"github.com/gorilla/websocket"
)

type Controllers struct {
	UserControllers.UserControllers
	RoomControllers.RoomControllers
}

func NewControllers(logger *logger.Logger, service *service.Service, upgrader *websocket.Upgrader) *Controllers {
	return &Controllers{
		UserControllers: UserControllers.NewUserControllers(*logger, service),
		RoomControllers: RoomControllers.NewRoomControllers(logger, service, upgrader),
	}
}
