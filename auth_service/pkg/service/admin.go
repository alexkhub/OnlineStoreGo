package service

import (
	authservice "auth_service"
	"auth_service/pkg/repository"
	"log"
	"net/url"

	"github.com/IBM/sarama"
)

type AdminService struct{
	repos repository.Admin
	jwt_service JWTManager
	producer sarama.SyncProducer

}

func NewAdminService(repos repository.Admin, jwt_service JWTManager, producer sarama.SyncProducer) *AdminService{
	return &AdminService{
		repos:  repos,
		jwt_service: jwt_service, 
		producer: producer,
		
	}
}


func (s *AdminService) UserList(filter url.Values)([]authservice.AdminUserListSerializer, error){
	return s.repos.UserListPostgres(filter)
}

func (s *AdminService) RoleList()([]authservice.RoleListSerializer, error){
	return s.repos.RoleListPostgres()
}

func (s *AdminService) UserBlock(user_id int)(error){
	
	err := s.repos.UserBlockPostgres(user_id)
	if err!= nil{
		return err 
	}
	data, err := s.repos.GetBlockDataPostgres(user_id)
	if err != nil{
		return err 
	}
	go func(){
		err = SendBlockKafkaMessage(s.producer, data)
		if err!= nil{
			log.Printf("Block Kafka %s", err.Error())
		}
	}()
	return nil
}
 
func (s *AdminService) UserUnblock(user_id int)(error){
	err :=  s.repos.UserUnblockPostgres(user_id)
	if err!= nil{
		return err 
	}
	data, err := s.repos.GetBlockDataPostgres(user_id)
	if err != nil{
		return err 
	}
	go func(){
		err = SendBlockKafkaMessage(s.producer, data)
		if err!= nil{
			log.Printf("Unblock Kafka %s", err.Error())
		}
	}()
	return nil
}