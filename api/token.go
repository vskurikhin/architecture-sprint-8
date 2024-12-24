package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Authorization = "Authorization"
	Bearer        = "Bearer "
)

var (
	ErrParsingValidatingToken = errors.New("error parsing or validating token")
	ErrInvalidToken           = errors.New("invalid token")
)

// GetToken получение строки токена из HTTP заголовка Authorization
func GetToken(request *http.Request) string {
	token := request.Header.Get(Authorization)
	splitToken := strings.Split(token, Bearer)
	if len(splitToken) < 2 {
		return ""
	}
	token = splitToken[1]
	return token
}

// ValidateToken валидация токена
func ValidateToken(reqToken string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, ErrParsingValidatingToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return token, nil
}
