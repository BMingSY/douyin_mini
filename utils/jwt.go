package utils

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/golang-jwt/jwt"
)

const (
	secretKey                = "DarkGreen"
	defaultExpireSeconds int = 36000
)

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

type MyClaims struct {
	UserID int `json:"userId"`
	jwt.StandardClaims
}

// 根据userID 生成 token
func GenerateToken(username string, userID int, expiredSeconds int) (tokenString string, err error) {
	if expiredSeconds == 0 {
		expiredSeconds = defaultExpireSeconds
	}
	mySigningKey := []byte(secretKey)
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	logs.Info("Token will be expired at ", time.Unix(expireAt, 0))

	claims := MyClaims{
		userID,
		jwt.StandardClaims{
			Issuer:    username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("error: failed to generate token")
	}
	return tokenStr, nil
}

// token 解析
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid
	}
}
