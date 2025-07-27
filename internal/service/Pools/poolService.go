package Pools

import (
	"context"
	"errors"
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	db "github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/redis"
	"github.com/SimpleOG/WebSocketChat/internal/service/Clients"
	"github.com/SimpleOG/WebSocketChat/internal/service/Rooms"
	"github.com/SimpleOG/WebSocketChat/util/config"
	"go.uber.org/zap"
	"sync"
)

// Проброс наружу методов
type Pools interface {
	ServePool(ctx context.Context)
	CheckRoom(ctx context.Context, hash string) error
	AddClientToRoom(ctx context.Context, hash string, clients *Clients.Client) error
	DeleteRoom(hash string)
}
type Pool struct {
	rooms      map[string]chan Clients.Client // кладём хеш комнаты и канал для записи новых пользователей
	newRoom    chan Rooms.Rooms
	deleteRoom chan string //сравниваем хеши и удаляет
	logger     logger.Logger
	mu         *sync.RWMutex
	querier    db.Querier
	redis      redis.RedisInterface
	config     config.Config
}

func NewPool(querier db.Querier, redisInterface redis.RedisInterface, config config.Config, logger logger.Logger) Pools {
	return &Pool{
		rooms:      make(map[string]chan Clients.Client),
		newRoom:    make(chan Rooms.Rooms, config.MaxEntryBuffSize),
		deleteRoom: make(chan string, config.MaxEntryBuffSize),
		logger:     logger,
		mu:         &sync.RWMutex{},
		querier:    querier,
		redis:      redisInterface,
		config:     config,
	}
}

// ServePool Запуск обслуживания всех комнат
func (p *Pool) ServePool(ctx context.Context) {
	p.logger.Info("Pools is serving")
	defer p.logger.Info("ServePool stopped")
	//Слушаем каналы для добавления и удаления комнат
	for {
		select {
		case _ = <-ctx.Done():
			return
		case room := <-p.newRoom:
			p.logger.Info("Пришла комната ", zap.String("№", room.GetRoomHash()))

			go room.ServeRoom(ctx)
		case hash := <-p.deleteRoom:
			go p.RoomDeleting(hash)

		}
	}
}

// Выставляем наружу метод, который проверяет комнату
func (p *Pool) CheckRoom(ctx context.Context, hash string) error {
	if _, ok := p.rooms[hash]; ok {
		return nil
	} else { //Если комнаты нет, то создаёт комнату
		err := p.CreateNewRoom(ctx, hash)
		if err != nil {
			return err
		}

	}
	return nil
}

func (p *Pool) CreateNewRoom(ctx context.Context, hash string) error {
	room := Rooms.NewRoom(hash, p.querier, p.redis, p.config, p.logger)
	err := p.querier.CreateRoom(ctx, hash)
	if err != nil {
		return err
	}
	p.mu.Lock()
	p.rooms[room.GetRoomHash()] = room.GetRoomChan()
	p.mu.Unlock()
	p.logger.Info("room  successfully added  ", zap.String("room", room.GetRoomHash()))
	if room == nil {
		p.logger.Error("Received nil room")
	}
	p.newRoom <- room
	return nil

}
func (p *Pool) DeleteRoom(hash string) {
	p.deleteRoom <- hash
}
func (p *Pool) RoomDeleting(hash string) {
	p.mu.Lock()
	for i := range p.rooms {
		if i == hash {
			p.logger.Info("room  successfully deleted  ", zap.String("room", hash))
			delete(p.rooms, hash)
		}
	}
	p.mu.Unlock()
}
func (p *Pool) AddClientToRoom(ctx context.Context, hash string, client *Clients.Client) error {
	room, ok := p.rooms[hash]
	if !ok {
		p.logger.Error("room is not found :", zap.String("roomID", hash))
		return errors.New(fmt.Sprintf("room %d is not found ", hash))
	}
	client.Redis = p.redis
	room <- *client
	return nil
}
