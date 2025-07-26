package UserControllers

import (
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/api/response"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserControllers interface {
	RegisterUser(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type UserController struct {
	logger  logger.Logger
	service *service.Service
}

func NewUserControllers(logger logger.Logger, service *service.Service) UserControllers {
	return &UserController{
		logger:  logger,
		service: service,
	}
}

func (u *UserController) RegisterUser(ctx *gin.Context) {
	var user db.CreateUserParams
	if err := ctx.ShouldBindJSON(&user); err != nil {
		u.logger.Error("cannot process users date", zap.Error(err))
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(err))
		return
	}
	_, err := u.service.AuthService.RegisterUser(ctx, user)
	if err != nil {
		u.logger.Error("cannot create new user :", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}
	u.logger.Info("User created sucessfully")
}
func (u *UserController) Login(ctx *gin.Context) {
	var userInfo db.GetUserForLoginParams
	if err := ctx.ShouldBindJSON(&userInfo); err != nil {
		u.logger.Error("cannot process users date", zap.Error(err))
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(err))
		return
	}
	userToken, err := u.service.AuthService.LoginUser(ctx, userInfo)
	if err != nil {
		u.logger.Error(fmt.Sprintf("cannot login  user %v : ", userInfo.Username), zap.Error(err))
		ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(err))
		return
	}
	u.logger.Info(fmt.Sprintf("User %d is authorized "))
	ctx.JSON(200, userToken)
}
