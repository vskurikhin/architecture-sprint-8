package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const PublicKey = "public_key"

var ErrNoPublicKey = errors.New("doen't public_key")

// GetKeycloakPublicKey в Keycloak RS256 используется по умолчанию, а публичный ключ, ф-ция его получения
func GetKeycloakPublicKey(keycloakURL string) (*rsa.PublicKey, error) {
	client := http.Client{}
	resp, errGet := client.Get(keycloakURL)
	if errGet != nil {
		return nil, errGet
	}
	bytes, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}
	issuer := make(map[string]interface{})
	errUnmarshal := json.Unmarshal(bytes, &issuer)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	var base64EncodedPublicKey string
	if b, ok0 := issuer[PublicKey]; !ok0 {
		return nil, ErrNoPublicKey
	} else {
		base64EncodedPublicKey = b.(string)
	}
	return parseKeycloakRSAPublicKey(base64EncodedPublicKey)
}

func parseKeycloakRSAPublicKey(base64Encoded string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}
