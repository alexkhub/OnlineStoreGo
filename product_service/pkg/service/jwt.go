package service

import (
	"encoding/json"
	"errors"
	"fmt"
	productservice "product_service"

	jwt "github.com/dgrijalva/jwt-go"
)



type Manager struct{
	signingKey string

}

func NewManager(signingKey string) *Manager{
	return &Manager{signingKey: signingKey}
}


func (m *Manager) Parse(accessToken string) (productservice.AuthMiddlewareSerializer, error) {
	var auth_parse productservice.AuthMiddlewareSerializer

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
		
		return auth_parse,  errors.New("error get user claims from token")
	}

	err = json.Unmarshal([]byte(claims["sub"].(string)), &auth_parse)
	if err != nil {
		fmt.Println(err)
		return auth_parse, err
	}
	
	return auth_parse, nil
}

