package service

import "order_service/pkg/repository"

type AdminService struct {
	repos repository.Admin
}

func NewAdminService(repos repository.Admin) *AdminService {
	return &AdminService{repos: repos}
}

func (s *AdminService) RemoveCartPoint(product_id int) error {
	return s.repos.RemoveCartPointPostgres(product_id)
}
