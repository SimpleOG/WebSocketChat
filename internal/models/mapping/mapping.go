package mapping

import (
	"encoding/json"
	"github.com/SimpleOG/WebSocketChat/internal/models"
	"github.com/go-redis/redis/v8"
)

func MapRedisMessageToClientMsg(redisMsg *redis.Message) (*models.ClientMessage, error) {
	var msg = new(models.ClientMessage)
	if err := json.Unmarshal([]byte(redisMsg.Payload), msg); err != nil {
		return nil, err
	}
	return msg, nil
}
