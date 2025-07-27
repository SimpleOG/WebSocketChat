package service

import (
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/redis"
	"github.com/SimpleOG/WebSocketChat/internal/service/Pools"
	"github.com/SimpleOG/WebSocketChat/internal/service/authService"
	auth "github.com/SimpleOG/WebSocketChat/pkg/JWTTokens"
	"github.com/SimpleOG/WebSocketChat/util/config"
)

type Service struct {
	AuthService AuthService.AuthorizationService
	Pool        Pools.Pools
}

func NewService(logger logger.Logger, maker auth.JWTMaker, querier db.Querier, redis redis.RedisInterface, config config.Config) *Service {
	return &Service{
		AuthService: AuthService.NewAuthService(maker, querier, logger),
		Pool:        Pools.NewPool(querier, redis, config, logger),
	}
}
