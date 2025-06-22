package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"time"

	"auth_service"

	jwt "github.com/dgrijalva/jwt-go"
)

type Manager struct {
	signingKey  string
	signingKey2 string
}

func NewManager(signingKey string, signingKey2 string) *Manager {
	return &Manager{signingKey: signingKey, signingKey2: signingKey}
}

func (m *Manager) CreateJwtAccess(user_id, role_id string) (string, error) {
	subject := make(map[string]string)
	subject["id"] = user_id
	subject["role"] = role_id
	subject_json, err := json.Marshal(subject)

	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		Subject:   string(subject_json),
	})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) CreateJwtRefresh(user_id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		Subject:   user_id,
	})
	return token.SignedString([]byte(m.signingKey2))

}

func (m *Manager) Parse(accessToken string) (authservice.AuthMiddlewareSerializer, error) {
	var auth_parse authservice.AuthMiddlewareSerializer

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return auth_parse, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {

		return auth_parse, errors.New("error get user claims from token")
	}

	err = json.Unmarshal([]byte(claims["sub"].(string)), &auth_parse)
	if err != nil {
		fmt.Println(err)
		return auth_parse, err
	}

	return auth_parse, nil
}

func (m *Manager) ParseRefreshToken(refreshToken string) (int, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey2), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {

		return 0, errors.New("error get user claims from token")
	}
	return strconv.Atoi(claims["sub"].(string))
}
