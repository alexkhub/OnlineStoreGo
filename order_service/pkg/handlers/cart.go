package handlers

import (
	"net/http"
	orderservice "order_service"

	"github.com/gin-gonic/gin"
	
)

func (h *Handler) MainPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "order_serive",
	})
}

func (h *Handler) GetMyCartHandler(c *gin.Context ){
	user, err := GetUserId(c)

	if err != nil{
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}
	data, err := h.services.Cart.CartList(user)
	if err != nil{
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListCartResponse{Data: data})
	
}

func (h *Handler) AddProductHandler(c *gin.Context){
	var input orderservice.CreateCartSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := GetUserId(c)

	if err != nil{
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}
	data, err := h.services.Cart.CreateCart(user, input.Product)
	if err != nil{
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	
	c.JSON(http.StatusOK, data)

}