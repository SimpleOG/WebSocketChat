package AuthService

import (
	"context"
	"errors"
	"fmt"
	"github.com/SimpleOG/WebSocketChat/internal/logger"
	"github.com/SimpleOG/WebSocketChat/internal/repositories/postgresql/sqlc"
	"github.com/SimpleOG/WebSocketChat/pkg/JWTTokens"
	"github.com/SimpleOG/WebSocketChat/util/hashing"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthorizationService interface {
	RegisterUser(ctx context.Context, userData db.CreateUserParams) (int32, error)
	GetUser(ctx context.Context, userId int32) (db.User, error)
	LoginUser(ctx context.Context, loginData db.GetUserForLoginParams) (string, error)
}

type Authorization struct {
	maker  JWTTokens.JWTMaker
	q      db.Querier
	logger logger.Logger
}

func NewAuthService(maker JWTTokens.JWTMaker, q db.Querier, logger logger.Logger) AuthorizationService {
	return &Authorization{
		maker:  maker,
		q:      q,
		logger: logger,
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
			return db.User{}, fmt.Errorf("no users was detected")
		}
	}

	return user, nil
}

func (a *Authorization) LoginUser(ctx context.Context, loginData db.GetUserForLoginParams) (string, error) {
	unameEmail := db.GetUserByUsernameOrEmailParams{
		Username: loginData.Username,
		Email:    loginData.Email,
	}

	//Сначала найти юзера и взять его хеш
	user, err := a.q.GetUserByUsernameOrEmail(ctx, unameEmail)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return "", err
	}
	claims := &jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(60 * time.Minute).Unix(),
	}
	userToken, err := a.maker.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return userToken, err
}
