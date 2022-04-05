package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	UserExpireDuration  = time.Hour * 2
	UserSecretKey       = "front_user_token"
	AdminExpireDuration = time.Hour * 48
	AdminSecretKey      = "admin_user_token"
)

type UserToken struct {
	jwt.StandardClaims
	Email  string `json:"username"`
	UserId string `json:"user_id"`
}

func GenerateToken(email string, expireTime time.Duration, secret string) (string, error) {
	//token=头部+载荷+签名
	//载荷=注册声明+私有声明
	user := UserToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireTime).Unix(),
			Issuer:    "yuyoung.web.sk_web",
		},
		Email: email,
	}

	//头部=算法+描述内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	//签名并生成
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString, nil
}

func AuthToken(tokenString string, secret string) (bool, string) {
	var user UserToken
	token, err := jwt.ParseWithClaims(tokenString, &user, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return false, ""
	}
	return true, user.Email
}
