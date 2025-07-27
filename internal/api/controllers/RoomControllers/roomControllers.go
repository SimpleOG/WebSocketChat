package RoomControllers

import (
	"errors"
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/api/response"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	"github.com/SimpleOG/WebSocketChat/internal/service/Clients"
	"github.com/SimpleOG/WebSocketChat/util/hashing"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type RoomControllers interface {
	Rooms(ctx *gin.Context)
	ServeRooms(ctx *gin.Context)
}
type Room struct {
	logger   *logger.Logger
	service  *service.Service
	upgrader websocket.Upgrader
}

func NewRoomControllers(logger *logger.Logger, service *service.Service, upgrader *websocket.Upgrader) RoomControllers {
	return &Room{
		logger:   logger,
		service:  service,
		upgrader: *upgrader,
	}
}

type RoomUsers struct {
	Users []int32 `json:"users"`
}

func GetIDParamFromCtx(ctx *gin.Context) int32 {

	ID := ctx.Value("id")
	if ID == nil {
		return 0
	}
	strVal, ok := ID.(string)
	if !ok {
		return 0
	}
	val, err := strconv.ParseInt(strVal, 10, 32)
	if err != nil {
		return 0
	}
	return int32(val)

}
func (r *Room) Rooms(ctx *gin.Context) {
	// запрос на эндпоинт сам юзер + те с кем он хочет общаться
	var users RoomUsers
	if err := ctx.ShouldBindJSON(&users); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(err))
		return
	}
	//описание логики
	/*
		Чел должен делать запрос и указывать с кем хочет общаться
		Система проверяет есть ли уже созданный для этих пользователей чат
		В конце концов юзер получает хеш чата для дальнейшего общения
	*/
	//берем id из токена
	userID := GetIDParamFromCtx(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(errors.New("cannot parse token")))
		return
	}
	users.Users = append(users.Users, userID)
	// создаем хеш и проверяем есть ли такой среди комнат если есть то создаем из юзера клиента
	hash := hashing.HashUsersForRoomUnique(users.Users)
	err := r.service.Pool.CheckRoom(ctx, hash)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"RoomHash": hash})

}

func (r *Room) ServeRooms(ctx *gin.Context) {
	a := r
	fmt.Println(a)
	//Получаем хеш комнаты из заголовка
	hash := ctx.GetHeader("hash")
	if hash == "" {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse(errors.New("no room hash provided")))
	}
	//Из токена берем id
	userID := GetIDParamFromCtx(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(errors.New("cannot parse token")))
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
		(*r.logger).Error("Не удалось улучшить соединение для пользователя  ", zap.Int32("id", user.ID))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}
	// Создаем клиента на основе данных юзера
	client := Clients.CreateClient(user, ws, r.logger)
	err = r.service.Pool.AddClientToRoom(ctx, hash, &client)
	if err != nil {
		(*r.logger).Error("Не удалось добавить юзера в комнату   ", zap.Int32("id", user.ID))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}
	go client.ReadMessageFromClient(ctx, hash)
	client.WriteMessageToClient(ctx)
}
