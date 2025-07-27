package Pools

import (
	"context"
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
	CheckRoom(ctx context.Context, hash string, client Clients.Clients)
	DeleteRoom(hash string)
}
type Pool struct {
	rooms      map[string]chan Clients.Clients // кладём хеш комнаты и канал для записи новых пользователей
	newRoom    chan Rooms.Rooms
	deleteRoom chan string //сравниваем хеши и удаляет
	logger     logger.Logger
	mu         sync.RWMutex
	querier    db.Querier
	redis      redis.RedisInterface
	config     config.Config
}

func NewPool(querier db.Querier, redisInterface redis.RedisInterface, config config.Config, logger logger.Logger) Pools {
	return &Pool{
		rooms:      make(map[string]chan Clients.Clients),
		newRoom:    make(chan Rooms.Rooms, config.MaxEntryBuffSize),
		deleteRoom: make(chan string, config.MaxEntryBuffSize),
		logger:     logger,
		mu:         sync.RWMutex{},
		querier:    querier,
		redis:      redisInterface,
		config:     config,
	}
}

// ServePool Запуск обслуживания всех комнат
func (p *Pool) ServePool(ctx context.Context) {
	//Слушаем каналы для добавления и удаления комнат
	for {
		select {
		case _ = <-ctx.Done():
			return
		case room := <-p.newRoom:
			go room.ServeRoom(ctx)
		case hash := <-p.deleteRoom:
			go p.RoomDeleting(hash)
		}
	}
}

// Выставляем наружу метод, который проверяет комнату
func (p *Pool) CheckRoom(ctx context.Context, hash string, client Clients.Clients) {
	//Если комната есть, то засовывает в неё нового юзера
	if room, ok := p.rooms[hash]; ok {
		room <- client
	} else { //Если комнаты нет, то создаёт комнату
		p.CreateNewRoom(ctx, hash)
	}
}
func (p *Pool) CreateNewRoom(ctx context.Context, hash string) {
	room := Rooms.NewRoom(hash, p.querier, p.redis, p.config)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rooms[room.GetRoomHash()] = room.GetRoomChan()
	p.logger.Info("room  successfully added  ", zap.String("room", room.GetRoomHash()))
	p.newRoom <- room
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
