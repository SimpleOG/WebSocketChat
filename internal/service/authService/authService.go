package AuthService

import (
	"context"
	"errors"
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	auth "github.com/SimpleOG/WebSocketChat/pkg/JWTTokens"
	"github.com/SimpleOG/WebSocketChat/util/hashing"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AuthorizationService interface {
	RegisterUser(ctx context.Context, userData db.CreateUserParams) (int32, error)
	GetUser(ctx context.Context, userId int32) (db.User, error)
	LoginUser(ctx context.Context, loginData db.GetUserForLoginParams) (db.User, error)
}

type Authorization struct {
	logger logger.Logger
	maker  auth.JWTMaker
	q      db.Querier
}

func NewAuthService(maker auth.JWTMaker, q db.Querier, logger logger.Logger) AuthorizationService {
	return &Authorization{
		maker:  maker,
		logger: logger,
		q:      q,
	}
}

// Регистрация нового пользователя
func (a *Authorization) RegisterUser(ctx context.Context, userData db.CreateUserParams) (int32, error) {
	HashedPass, err := hashing.GeneratePassword(userData.Password)
	if err != nil {
		return 0, err
	}
	userData.Password = HashedPass
	NewUser, err := a.q.CreateUser(ctx, userData)
	if err != nil {
		return 0, fmt.Errorf("failed to create  new user %v", err)
	}
	return NewUser.ID, nil
}
func (a *Authorization) GetUser(ctx context.Context, userId int32) (db.User, error) {
	user, err := a.q.GetUsersById(ctx, userId)
	if err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			a.logger.Error("Отсуствует пользователь с id ", zap.Int32("id", user.ID))
			return db.User{}, fmt.Errorf("no users was detected")
		}
	}

	return user, nil
}

func (a *Authorization) LoginUser(ctx context.Context, loginData db.GetUserForLoginParams) (db.User, error) {

	hashedPass, err := hashing.GeneratePassword(loginData.Password)
	if err != nil {
		return db.User{}, err
	}
	loginData.Password = hashedPass
	User, err := a.q.GetUserForLogin(ctx, loginData)
	if err != nil {
		return db.User{}, err
	}
	return User, err
}
