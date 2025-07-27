package main

import (
	"context"
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/api/middlewares"
	"github.com/SimpleOG/WebSocketChat/internal/api/server"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/redis"
	"github.com/SimpleOG/WebSocketChat/internal/service"
	auth "github.com/SimpleOG/WebSocketChat/pkg/JWTTokens"
	"github.com/SimpleOG/WebSocketChat/util/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
)

func main() {

	config, err := config.InitConfig("./")
	if err != nil {
		log.Fatalf("error while downloading config :%v", err)
	}
	logger, err := logger.NewLogger(zapcore.Level(config.LoggerLevel))
	ctx := context.TODO()
	connPool, err := pgxpool.New(ctx, config.DBDSource)
	if err != nil {
		logger.Fatal("Cannot create  connection pool :", zap.Error(err))
		return
	}
	err = connPool.Ping(ctx)
	if err != nil {
		logger.Fatal("Cannot ping db :", zap.Error(err))
		return
	}
	logger.Info("Db successfuly started ")
	querier := db.New(connPool)
	if err := runDBMigration(config.MigrationUrl, config.DBDSource); err != nil && err.Error() != "no change" {
		logger.Fatal("Cannot start migration", zap.Error(err))
		return
	}
	redis, err := redis.NewRedisClient(ctx, logger, config)
	if err != nil {
		logger.Fatal("Cannot connect to redis :", zap.Error(err))
		return
	}
	jwtMaker := auth.NewJWTMaker(config.SecretKey)

	service := service.NewService(logger, jwtMaker, querier, redis, config)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	middlewares := middlewares.NewMiddleware(jwtMaker)
	server := server.NewServer(&logger, router, service, middlewares, &upgrader)

	if err := server.Run(ctx, config.ServerAddress); err != nil {
		logger.Fatal("Cannot start server  :", zap.Error(err))
		return
	}

	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	//<-quit
	////logrus.Print("TodoApp server shutting")
	////if err := srv.Shutdown(context.Background()); err != nil {
	////	logrus.Errorf("ну всё кранты из за того что %s", err.Error())
	////}
}
func runDBMigration(migrationURL, dbSource string) error {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return fmt.Errorf("cannot find migration to up %v", err)
	}
	if err = migration.Up(); err != nil {
		if err.Error() == "no change" {
			return err
		}
		return fmt.Errorf("cannot start migration %v", err)
	}
	return nil
}
func StopDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("cannot find migration to down", err)
	}
	if err = migration.Down(); err != nil {
		log.Fatalln("cannot stop migration", err)
	}
	log.Println("stopped")
}
