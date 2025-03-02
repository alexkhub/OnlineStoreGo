package service

import(
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