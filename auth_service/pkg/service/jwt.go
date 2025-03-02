package service

import (
	"encoding/json"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)


type Manager struct{
	signingKey string
	signingKey2 string
}

func NewManager(signingKey string, signingKey2 string) *Manager{
	return &Manager{signingKey: signingKey, signingKey2: signingKey}
}

func (m *Manager) CreateJwtAccess(user_id, role_id string) (string, error){
	subject := make(map[string]string)
	subject["id"] = user_id
	subject["role"] = role_id
	subject_json, err := json.Marshal(subject)

	if err != nil{
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		Subject: string(subject_json),
	})
	return token.SignedString([]byte(m.signingKey))
}


func (m *Manager) CreateJwtRefresh(user_id string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		Subject: user_id,
	})
	return token.SignedString([]byte(m.signingKey2))

}