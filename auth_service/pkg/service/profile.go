package service

import (
	authservice "auth_service"
	"auth_service/pkg/repository"
)

type ProfileService struct{
	repos repository.Profile

}

func NewProfileService(repos repository.Profile ) *ProfileService{
	return &ProfileService{
		repos:  repos,
	}
}

func (s *ProfileService) UserProfile( user_id int) (authservice.ProfileSerializer, error){
	
	return s.repos.UserProfilePostgres(user_id)
}