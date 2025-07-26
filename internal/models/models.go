package models

import (
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"golang.org/x/net/websocket"
)

type ChatClient struct {
	User db.User
	Conn *websocket.Conn
}
type ClientMessage struct {
	User_id     int32
	Msg_content string
}
