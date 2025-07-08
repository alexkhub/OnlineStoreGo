package service

import (
	"context"
	productservice "product_service"
	"product_service/pkg/repository"
	grpc_product_service "github.com/alexkhub/OnlineStoreProto/gen/go/comment"
	"github.com/redis/go-redis/v9"
)

type CommentService struct {
	repos   repository.Comment
	redisDB *redis.Client
	gRPCComment grpc_product_service.CommentClient
}

func NewCommentService(repos repository.Comment, redisDB *redis.Client, gRPCComment grpc_product_service.CommentClient) *CommentService {
	return &CommentService{repos: repos, redisDB: redisDB, gRPCComment: gRPCComment}
}

func (s *CommentService) CreateComment(data productservice.CreateCommentSerializer, product_id, user_id int) (int, error) {
	return s.repos.CreateCommentPostgres(data, product_id, user_id)
}

func (s *CommentService) RemoveUserComment(user_id int) error {
	return s.repos.RemoveUserCommentPostgres(user_id)
}

func (s *CommentService) CommentList(product_id int) ([]productservice.ListCommentSerializer, error){
	var commentList []productservice.ListCommentSerializer 

	commentData, err := s.repos.CommentListPostgres(product_id)
	if err != nil{
		return nil, err
	}

	if len(commentData) == 0{
		return commentList, nil
	}
	UniqueUserIdMap := make(map[int64]struct{},) 

	for _, value :=  range commentData{
		UniqueUserIdMap[value.Id] = struct{}{}
	}

	UniqueUserId := make([]int64, 0, len(UniqueUserIdMap))

	for key := range UniqueUserIdMap{
		UniqueUserId = append(UniqueUserId, key)
	}
	userData, err := s.gRPCComment.GetUserData(context.Background(), &grpc_product_service.CommentIdRequest{Id: UniqueUserId})
	if err != nil{
		return nil, err
	}

	userDataMap := make(map[int64]*grpc_product_service.UserData )
	for _, value := range userData.Data{
		userDataMap[value.Id] = value
	}

	for _, value := range commentData{
		user := userDataMap[int64(value.User)]
		commentList = append(commentList, productservice.ListCommentSerializer{
			Id: value.Id,
			Title: value.Title,
			Raiting: value.Raiting,
			User:  productservice.ComentUserDataSerializer{
				Id: user.Id,
				FullName: user.FullName,
				Image: user.Image,
			},
			Message: value.Message,
			CreateAt: value.CreateAt,
		})
	}
	return commentList, nil
}


func (s *CommentService) RemoveComment(comment_id int, user_id int) error{
	return s.repos.RemoveCommentPostgres(comment_id, user_id)
}