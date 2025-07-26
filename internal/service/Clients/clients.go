package Clients

import (
	"context"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/models"
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/redis"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	UserInfo db.User
	MsgChan  chan string // канал в который приходят новые сообщения
	conn     *websocket.Conn
	logger   logger.Logger
	redis    redis.RedisInterface
}

// Считываем всё что клиент пишет в соединение вебсокета
func (c *Client) ReadMessageFromClient(ctx context.Context, roomHash string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.logger.Info("WebSocket закрыт клиентом")
					return
				}
				c.logger.Error("error while reading user message",
					zap.Error(err),
					zap.Int32("user: ", c.UserInfo.ID))
				continue
			}
			message := models.ClientMessage{
				User_id:     c.UserInfo.ID,
				Msg_content: string(msg),
			}
			//Кладём сообщение в канал редиса
			err = c.redis.SendMessageToChan(ctx, roomHash, message)
			if err != nil {
				c.logger.Error("error while trying to send to redis  message",
					zap.Error(err),
					zap.Int32("user: ", c.UserInfo.ID))
				continue

			}
		}
	}
}

// Отправляем в вебсокет сообщения
func (c *Client) WriteMessageToClient(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Writing to client connection is stopped")
			return
		case msg := <-c.MsgChan:
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				c.logger.Error("error while sending msg to user's connection", zap.Error(err))
				continue
			}
		}
	}
}
