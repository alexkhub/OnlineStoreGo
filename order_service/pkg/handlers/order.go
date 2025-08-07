package handlers

import (
	"net/http"
	orderservice "order_service"

	"github.com/gin-gonic/gin"
	v "github.com/asaskevich/govalidator"
)


func (h *Handler) PaymentMethodeListHandler(c *gin.Context){
	data, err := h.services.Order.PaymentMethodeList()
	if err != nil{
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListPaymentMethodeResponse{Data: data})
}


func (h *Handler) CreateOrderHandler(c *gin.Context){
	var input orderservice.CreateOrderSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := GetUserId(c)

	if err != nil{
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	input.User = user

	id, err := h.services.Order.CreateOrder(input)

	if err != nil{
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id" : id,
	})

}