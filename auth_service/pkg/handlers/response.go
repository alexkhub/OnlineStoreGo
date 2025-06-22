package handlers

import (
	authservice "auth_service"
	"log"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getListUserResponse struct {
	Data []authservice.AdminUserListSerializer `json:"data"`
}

type getListRoleResponse struct {
	Data []authservice.RoleListSerializer `json:"data"`
}

func newErrorMessage(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
