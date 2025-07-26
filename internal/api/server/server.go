package server

import (
	"context"
	"github.com/SimpleOG/WebSocketChat/internal/api/controllers"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	"github.com/SimpleOG/WebSocketChat/internal/service/Pools"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	controllers *controllers.Controllers
	pools       Pools.Pools
}

func NewServer(logger logger.Logger, engine *gin.Engine, service *service.Service, pools Pools.Pools) *Server {
	return &Server{
		router:      engine,
		controllers: controllers.NewControllers(logger, service),
		pools:       pools,
	}
}
func (s *Server) Run(ctx context.Context, addr string) error {
	s.pools.ServePool(ctx)
	s.SetupRoutes()
	return s.router.Run(addr)
}
func (s *Server) SetupRoutes() {
	api := s.router.Group("/api")
	s.SetupRoomsRoutes(api)
	s.SetupAuthRoutes(api)
}
func (s *Server) SetupRoomsRoutes(api *gin.RouterGroup) {
	rooms := api.Group("/rooms")
	{
		rooms.POST("/room", s.controllers.RoomControllers.ServeRooms)
	}
}
func (s *Server) SetupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", s.controllers.UserControllers.RegisterUser)
		auth.POST("/login", s.controllers.UserControllers.Login)
	}
}
func (s *Server) Setup3Routes() {}
