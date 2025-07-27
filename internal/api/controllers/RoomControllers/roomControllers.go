package RoomControllers

import (
	"errors"
	"github.com/SimpleOG/WebSocketChat/internal/api/response"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	"github.com/SimpleOG/WebSocketChat/internal/service/Clients"
	"github.com/SimpleOG/WebSocketChat/util/hashing"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type RoomControllers interface {
	ServeRooms(ctx *gin.Context)
}
type Room struct {
	logger   logger.Logger
	service  *service.Service
	upgrader websocket.Upgrader
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
		return
	}

	//описание логики
	/*
		Чел должен делать запрос и указывать с кем хочет общаться
		Система проверяет есть ли уже созданный для этих пользователей чат
		Если чат есть , добавляет этого юзера в чат
	*/
	//берем id из токена
	ID, ok := ctx.Get("id")
	if !ok {
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(errors.New("пользователь не найден")))
		r.logger.Error("отсутствует id текущего пользователя")
		return
	}
	userID, ok := ID.(int32)
	if !ok {
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(errors.New("токен содержит неправильные данные")))
		r.logger.Error("нет подходящего формата id внутри токена", zap.Any("id", userID))
		return
	}
	//Получаем инфу про текущего юзера из бд

	user, err := r.service.AuthService.GetUser(ctx, userID)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}
	//Если юзер есть, то начинается процесс открытия ws соединения
	ws, err := r.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		r.logger.Error("Не удалось улучшить соединение для пользователя  ", zap.Int32("id", user.ID))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}
	// создаем хеш и проверяем есть ли такой среди комнат если есть то создаем из юзера клиента
	hash := hashing.HashUsersForRoomUnique(users.users)
	client := Clients.CreateClient(user, ws, r.logger)
	r.service.Pool.CheckRoom(ctx, hash, &client)
	go client.ReadMessageFromClient(ctx, hash)
	client.WriteMessageToClient(ctx)

}
