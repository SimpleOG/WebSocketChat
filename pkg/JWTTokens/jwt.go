package JWTTokens

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

type JWTMaker interface {
	CreateToken(claims *jwt.MapClaims) (string, error)
	VerifyToken(tokenStr string) (jwt.MapClaims, error)
}

type Maker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) JWTMaker {
	return &Maker{secretKey: secretKey}
}

func (j *Maker) CreateToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("cannot sign token %v", err)
	}
	return tokenStr, nil
}
func (j *Maker) VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while pasring token %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
