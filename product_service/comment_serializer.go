package productservice

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type CreateCommentSerializer struct {
	Title   string `json:"title" binding:"required" valid:"-"`
	Message string `json:"message" valid:"-"`
	Raiting int    `json:"raiting"  binding:"required" valid:"-"`
}

type ComentUserDataSerializer struct {
	Id            int64 `json:"id" `
	FullName      string `json:"full_name"`
	Image         string `json:"image"` 
	
}

type ListCommentPostgresSerializer struct {
	Id int64 `json:"id" db:"id"`
	Title string  `json:"title" db:"title"`
	Raiting int    `json:"raiting" db:"raiting"`
	User int `json:"user" db:"user_id"`
	Message null.String `json:"message" db:"message"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}

type ListCommentSerializer struct {
	Id int64 `json:"id"`
	Title string  `json:"title" `
	Raiting int    `json:"raiting"`
	User ComentUserDataSerializer `json:"user"`
	CreateAt time.Time `json:"create_at" `
	Message null.String `json:"message" db:"message"`
}

