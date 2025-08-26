package handlers

import (
	"log"
	orderservice "order_service"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type MyError struct {
	Error string
	Code int64
}

func newErrorMessage(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

type getListCartResponse struct {
	Data []orderservice.CartSerializer `json:"data"`
}
type getListPaymentMethodeResponse struct {
	Data []orderservice.PaymentMethodeSerializer `json:"data"`
}

type getListUserOrderListResponse struct{
	Data []orderservice.UserOrderListSerializer `json:"data"`
}