package productservice

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type CreateCommentSerializer struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message"`
	Rating int    `json:"rating"  binding:"required"`
}

type CommentUserDataSerializer struct {
	Id            int64 `json:"id" `
	FullName      string `json:"full_name"`
	Image         string `json:"image"` 
	
}

type ListCommentPostgresSerializer struct {
	Id int64 `json:"id" db:"id"`
	Title string  `json:"title" db:"title"`
	Rating int    `json:"raiting" db:"rating"`
	User int `json:"user" db:"user_id"`
	Message null.String `json:"message" db:"message"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}

type ListCommentSerializer struct {
	Id int64 `json:"id"`
	Title string  `json:"title" `
	Rating int    `json:"raiting"`
	User CommentUserDataSerializer `json:"user"`
	CreateAt time.Time `json:"create_at" `
	Message null.String `json:"message" db:"message"`
}

