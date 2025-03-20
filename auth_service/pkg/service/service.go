package service

import (
	"auth_service/pkg/repository"
	"auth_service"
	"github.com/IBM/sarama"
	
)

const (
	AuthTopic = "auth_topic"
	ConfirmTopic = "confirm_topic"
)

type Authorization interface{
	Registration(user authservice.AuthRegistrationSerializer) (authservice.AuthRegistrationResponseSerializer, error)
	ActivateUser(id int)(error)
	LoginUser(user authservice.LoginUser) (authservice.JWTToken, error)
	RefreshJWTToken(refresh string) (authservice.JWTToken, error)
}

type Profile interface{
	UserProfile( user_id int) (authservice.ProfileSerializer, error)
}

type JWTManager interface{   
	CreateJwtAccess(user_id, role_id string) (string, error)
	CreateJwtRefresh(user_id string) (string, error)
	Parse(accessToken string) (authservice.AuthMiddlewareSerializer, error)
	ParseRefreshToken(refreshToken string) (int, error)
}

type Deps struct {
    Repos *repository.Repository
	JWTManager JWTManager
	Producer sarama.SyncProducer

}

type Service struct {
	Authorization
	Profile
}

func NewService(deps Deps) *Service{
    new_auth_service := NewAuthService(deps.Repos.Authorization, deps.JWTManager, deps.Producer)
    new_profile_service := NewProfileService(deps.Repos.Profile)
	
	return &Service{
		Authorization: new_auth_service,
		Profile: new_profile_service,
	}
}