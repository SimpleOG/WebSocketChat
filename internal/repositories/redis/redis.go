package redis

import (
	"context"
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/util/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisInterface interface {
	SendMessageToChan(ctx context.Context, Chan string, msg any) error
	SubOnChannel(ctx context.Context, chanName string) *redis.PubSub
}

type Redis struct {
	logger logger.Logger
	redis  *redis.Client
}

func NewRedisClient(ctx context.Context, logger logger.Logger, config config.Config) (RedisInterface, error) {
	client := redis.NewClient(&redis.Options{Addr: config.RedisAddress, Password: config.RedisPassword, DB: config.RedisDB})
	logger.Info("connecting to redis ", zap.String("host", config.RedisAddress), zap.Int32("db", int32(config.RedisDB)))
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot ping redis client : %v ", err)
	}
	RedisClient := &Redis{
		logger: logger,
		redis:  client,
	}

	return RedisClient, nil
}
func (r *Redis) SendMessageToChan(ctx context.Context, Chan string, msg any) error {
	r.logger.Info("Sending msg to chanel : ", zap.String("chan", Chan))
	err := r.redis.Publish(ctx, Chan, msg).Err()
	if err != nil {
		return fmt.Errorf("cannot publish msg to channel : %v", err)
	}
	return nil
}
func (r *Redis) SubOnChannel(ctx context.Context, chanName string) *redis.PubSub {
	return r.redis.Subscribe(ctx, chanName)
}
