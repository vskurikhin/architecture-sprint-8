package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"time"

	c0env "github.com/caarlos0/env"
)

const RetryCount = 100

type environments struct {
	KeycloakURL   string `env:"KEYCLOAK_URL"`
	KeycloakRealm string `env:"KEYCLOAK_REALM"`
}

var (
	keycloakURL = "http://localhost:8080/realms/reports-realm"
)

func init() {
	env, err := getEnvironments()
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	if env.KeycloakURL != "" {
		keycloakURL = env.KeycloakURL
	}
	if env.KeycloakRealm != "" {
		keycloakURL = fmt.Sprintf("%s/realms/%s", keycloakURL, env.KeycloakRealm)
	}
}

func main() {
	var publicKey *rsa.PublicKey
	errRetry := Retry(RetryCount, RetryCount*time.Microsecond, func() (err error) {
		publicKey, err = GetKeycloakPublicKey(keycloakURL)
		return err
	})
	if errRetry != nil {
		log.Fatal(errRetry)
	}
	http.HandleFunc("/reports", Handler(publicKey))
	log.Println(http.ListenAndServe(":8000", nil))
}

func getEnvironments() (*environments, error) {
	env := new(environments)
	err := c0env.Parse(env)
	return env, err
}
