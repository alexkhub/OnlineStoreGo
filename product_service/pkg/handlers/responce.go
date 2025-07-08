package handlers

import (
	"log"
	productservice "product_service"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getListCategoryResponse struct {
	Data []productservice.CategorySerializer `json:"data"`
}

type getListProductResponse struct {
	Data []productservice.ProductListSerailizer `json:"data"`
}

type getListCommentResponse struct {
	Data []productservice.ListCommentSerializer `json:"data"`
}

func newErrorMessage(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
