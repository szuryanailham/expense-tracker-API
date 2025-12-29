package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/szuryanailham/expense-tracker/internal/env"
)

type JWTToken struct {
	secretkey string
}

func CreateJWT(secret []byte, userID string) (string, error) {
	expirationSeconds := env.GetInt("JWT_EXPIRATION_IN_SECOND", 3600)
	expiration := time.Second * time.Duration(expirationSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id": userID,
		"expiredAt":time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenStr string, secret []byte)(string, error) {
	token, err := jwt.Parse(tokenStr,func(t *jwt.Token)(interface{},error){
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims , ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	userID,ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found")
	}
	exp, ok := claims["expiredAt"].(float64)
	if !ok {
		return "", errors.New("exp claim not found")
	}
	if int64(exp) < time.Now().Unix() {
		return "", errors.New("token expired")
	}
	return userID, nil
}

