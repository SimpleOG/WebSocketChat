package Rooms

import (
	"context"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/models"
	"github.com/SimpleOG/WebSocketChat/internal/models/mapping"
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/redis"
	"github.com/SimpleOG/WebSocketChat/internal/service/Clients"
	"github.com/SimpleOG/WebSocketChat/util/config"
	"go.uber.org/zap"
	"sync"
)

//Логика взаимодействия в рамках одной комнаты, где могут быть собраны несколько пользователей

type Rooms interface {
	GetRoomHash() string
	ServeRoom(ctx context.Context)
	GetRoomChan() chan Clients.Clients
}

type Room struct {
	roomHash         string //Хэш комнаты (состоит из хеша всех id)
	redis            redis.RedisInterface
	querier          db.Querier
	newClient        chan Clients.Clients          //Для добавления новых пользователей в чат
	clientsForDelete chan int32                    //Для удаления пользователей из чата
	msgChan          chan models.ClientMessage     //Для сообщений
	clients          map[*Clients.Clients]struct{} // Текущие участники чата
	mu               sync.RWMutex
	logger           logger.Logger
}

func NewRoom(roomHash string, querier db.Querier, redis redis.RedisInterface, config config.Config) Rooms {
	return &Room{
		roomHash:         roomHash,
		redis:            redis,
		querier:          querier,
		msgChan:          make(chan models.ClientMessage, config.MaxMsgBuffSize),
		clientsForDelete: make(chan int32, config.MaxEntryBuffSize),
		newClient:        make(chan Clients.Clients, config.MaxEntryBuffSize),
		clients:          make(map[*Clients.Clients]struct{}, config.MaxEntryBuffSize),
	}
}

// Для проброса хеша наверх
func (r *Room) GetRoomHash() string {
	return r.roomHash
}
func (r *Room) GetRoomChan() chan Clients.Clients {
	return r.newClient
}
func (r *Room) ReadRedisChannel(ctx context.Context) {
	sub := r.redis.SubOnChannel(ctx, r.roomHash)
	defer sub.Close()
	redisCh := sub.Channel()
	for {
		select {
		case msg := <-redisCh:
			clientMsg, err := mapping.MapRedisMessageToClientMsg(msg)
			if err != nil {
				r.logger.Error("error while mapping message from redis", zap.Error(err))
				continue
			}
			msgParams := db.CreateMessageParams{
				RoomID:   r.roomHash,
				SenderID: clientMsg.User_id,
				Content:  clientMsg.Msg_content,
			}
			err = r.querier.CreateMessage(ctx, msgParams)
			if err != nil {
				r.logger.Error("error while saving msg to database", zap.Error(err))
				continue
			}
			r.msgChan <- *clientMsg
		}

	}
}

// Функция для обслуживания комнаты
func (r *Room) ServeRoom(ctx context.Context) {
	r.logger.Info("room started ", zap.String("room", r.roomHash))
	go r.ReadRedisChannel(ctx)
	for {
		select {
		case newClient := <-r.newClient:
			go r.AddNewClientIntoChat(&newClient)
		case ClientForDelete := <-r.clientsForDelete:
			go r.DeleteClientFromRoom(ClientForDelete)
		case msg := <-r.msgChan:
			go r.ProcessMessage(msg)
		}

	}
}
func (r *Room) AddNewClientIntoChat(client *Clients.Clients) {
	r.mu.Lock()
	r.clients[client] = struct{}{}
	r.logger.Info("new user added into room ")
	r.mu.Unlock()
}
func (r *Room) DeleteClientFromRoom(id int32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, _ := range r.clients {
		if i.UserInfo.ID == id {
			delete(r.clients, i)
			r.logger.Info("user deleted from room ")
			break
		}
	}

}
func (r *Room) ProcessMessage(msg models.ClientMessage) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	//Отправляем смску всем юзерам внутри комнаты
	for i := range r.clients {
		go func() {
			i.MsgChan <- msg.Msg_content
		}()
	}
}
