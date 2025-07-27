package middlewares

import (
	"fmt"
	"github.com/SimpleOG/WebSocketChat/pkg/JWTTokens"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Middleware interface {
	ValidateToken(ctx *gin.Context)
}

type MiddlewareStruct struct {
	jwtMaker JWTTokens.JWTMaker
}

func NewMiddleware(maker JWTTokens.JWTMaker) Middleware {
	return &MiddlewareStruct{jwtMaker: maker}
}
func (m *MiddlewareStruct) ValidateToken(ctx *gin.Context) {
	//Чек что токен в принципе есть
	Bearer := ctx.GetHeader("Authorization")
	if Bearer == "" {
		ctx.AbortWithStatus(401)
	}
	token := strings.Split(Bearer, " ")[1]
	// проверяем правильность токена
	userClaims, err := m.jwtMaker.VerifyToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err})
	}
	user_id := fmt.Sprintf("%v", (userClaims)["sub"])
	ctx.Set("id", user_id)
	ctx.Next()

}
