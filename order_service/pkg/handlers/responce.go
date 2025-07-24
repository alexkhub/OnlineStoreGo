package handlers

import (
	"log"
	orderservice "order_service"

	"github.com/gin-gonic/gin"
)


type errorResponse struct {
	Message string `json:"message"`
}


func newErrorMessage(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}


type getListCartResponse struct {
	Data []orderservice.CartSerializer `json:"data"`
}