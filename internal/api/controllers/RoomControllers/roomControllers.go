package RoomControllers

import (
	"github.com/SimpleOG/WebSocketChat/internal/api/response"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	"github.com/SimpleOG/WebSocketChat/internal/service/Clients"
	"github.com/SimpleOG/WebSocketChat/util/hashing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoomControllers interface {
	ServeRooms(ctx *gin.Context)
}
type Room struct {
	logger  logger.Logger
	service *service.Service
}

func NewRoomControllers(logger logger.Logger, service *service.Service) RoomControllers {
	return &Room{
		logger:  logger,
		service: service,
	}
}

type RoomUsers struct {
	users []int32
}

func (r *Room) ServeRooms(ctx *gin.Context) {
	// запрос на эндпоинт сам юзер + те с кем он хочет общаться
	var users RoomUsers
	if err := ctx.ShouldBindJSON(users); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(err))
	}
	//Создаем из юзера клиента
	currentClient := Clients.Client{}
	// создаем хеш и проверяем есть ли такой среди комнат
	hash := hashing.HashUsersForRoomUnique(users.users)
	//если есть то присоединяем юзера к комнате
	r.service.Pool.CheckRoom(ctx, hash, currentClient)
	//если нет, то создаем комнату и закидываем в бд

}
