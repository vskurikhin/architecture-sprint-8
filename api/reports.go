package main

import (
	"crypto/rsa"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const ProtheticUser = "prothetic_user"

var (
	headers = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, Token",
	}
)

// Handler возвращает ф-цию обработчик для /reports
func Handler(publicKey *rsa.PublicKey) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(writer)
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}
		reqToken := GetToken(request)
		jwtToken, errValidateToken := ValidateToken(reqToken, publicKey)
		if errValidateToken != nil {
			log.Printf("error: %v\n", errValidateToken)
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write([]byte(errValidateToken.Error()))
			return
		}
		roles := RolesClaims{
			Claims: jwtToken.Claims.(jwt.MapClaims),
		}.ToMap().
			ToSlice().
			ToRoles()
		if _, ok := roles[ProtheticUser]; !ok {
			log.Printf("error: Forbidden\n")
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("Ok"))
	}
}

func enableCors(w http.ResponseWriter) {
	for header, value := range headers {
		w.Header().Set(header, value)
	}
}
