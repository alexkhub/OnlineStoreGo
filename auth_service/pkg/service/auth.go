package service

import (
	"auth_service"
	"auth_service/pkg/repository"
	"errors"
	"log"
	"strconv"

	"github.com/IBM/sarama"
)

type AuthService struct {
	repos       repository.Authorization
	jwt_service JWTManager
	producer    sarama.SyncProducer
}

func NewAuthService(repos repository.Authorization, jwt_service JWTManager, producer sarama.SyncProducer) *AuthService {
	return &AuthService{
		repos:       repos,
		jwt_service: jwt_service,
		producer:    producer,
	}
}

func (s *AuthService) Registration(user authservice.AuthRegistrationSerializer) (authservice.AuthRegistrationResponseSerializer, error) {
	if user.Password != user.RepeatPassword {
		return authservice.AuthRegistrationResponseSerializer{}, errors.New("the password will not match")
	}
	user.Password, _ = HashPassword(user.Password)
	id, email, err := s.repos.RegistrationPostrgres(user)
	if err != nil {
		return authservice.AuthRegistrationResponseSerializer{}, err
	}
	err = SendVerifyKafkaMessage(s.producer, authservice.AuthRegistrationResponseSerializer{Id: id, Email: email})

	if err != nil {
		log.Printf("Registration Kafka %s", err.Error())
	}
	return authservice.AuthRegistrationResponseSerializer{Id: id, Email: email}, nil
}

func (s *AuthService) ActivateUser(id int) error {
	return s.repos.ActivateUserPostgres(id)
}

func (s *AuthService) LoginUser(user authservice.LoginUser) (authservice.JWTToken, error) {
	var param string
	var value string

	if user.Username != "" {
		param = "username"
		value = user.Username

	} else if user.Email != "" {
		param = "email"
		value = user.Email
	} else {
		return authservice.JWTToken{}, errors.New("username and email is empty")
	}

	data, err := s.repos.LoginPostgres(param, value)
	if err != nil {
		return authservice.JWTToken{}, err
	}

	if !data.Activate || data.Block {
		return authservice.JWTToken{}, errors.New("no access to account ")
	}

	if !CheckPasswordHash(user.Password, data.Password) {
		return authservice.JWTToken{}, errors.New("password error")
	}

	access, err := s.jwt_service.CreateJwtAccess(strconv.Itoa(data.Id), strconv.Itoa(data.Role))
	if err != nil {
		return authservice.JWTToken{}, errors.New("JWT error " + err.Error())
	}

	refresh, err := s.jwt_service.CreateJwtRefresh(strconv.Itoa(data.Id))
	if err != nil {
		return authservice.JWTToken{}, errors.New("JWT error " + err.Error())
	}

	err = s.repos.CreateJwtRefreshPostgres(data.Id, refresh)
	if err != nil {
		log.Printf("refresh token not recorded in the DB = %d, error = %s", data.Id, err)
	}

	return authservice.JWTToken{Access: access, Refresh: refresh}, nil
}

func (s *AuthService) RefreshJWTToken(refresh string) (authservice.JWTToken, error) {

	user_id, err := s.jwt_service.ParseRefreshToken(refresh)
	if err != nil {
		return authservice.JWTToken{}, errors.New("refresh error " + err.Error())
	}

	data, err := s.repos.RefreshCheckUserPostgres(user_id)

	if err != nil {
		return authservice.JWTToken{}, errors.New("refresh error " + err.Error())
	}

	if !data.Activate || data.Block {
		return authservice.JWTToken{}, errors.New("no access to account ")
	}

	new_access, err := s.jwt_service.CreateJwtAccess(strconv.Itoa(user_id), strconv.Itoa(data.Role))
	if err != nil {
		return authservice.JWTToken{}, errors.New("JWT error " + err.Error())
	}

	new_refresh, err := s.jwt_service.CreateJwtRefresh(strconv.Itoa(user_id))
	if err != nil {
		return authservice.JWTToken{}, errors.New("JWT error " + err.Error())
	}

	err = s.repos.UpdateJwtRefreshPostres(user_id, refresh, new_refresh)
	if err != nil {
		return authservice.JWTToken{}, err
	}

	return authservice.JWTToken{Access: new_access, Refresh: new_refresh}, nil
}

func (s *AuthService) DeleteRefreshJWTToken(refresh string) error {
	return s.repos.DeleteRefreshJWTTokenPostgres(refresh)
}

func (s *AuthService) CloseAllSessions(id int) error {
	return s.repos.CloseAllSessionsPostgres(id)
}
