package middlewares

import (
	"github.com/SimpleOG/WebSocketChat/pkg/JWTTokens"
	"github.com/gin-gonic/gin"
	"net/http"
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
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatus(401)
	}
	// проверяем правильность токена
	userClaims, err := m.jwtMaker.VerifyToken(token)
	user_id := (*userClaims)["sub"].(string)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err})
	}
	ctx.Set("id", user_id)
	ctx.Next()

}
