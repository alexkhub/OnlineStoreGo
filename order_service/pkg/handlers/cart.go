package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	orderservice "order_service"
	"strconv"
)

func (h *Handler) MainPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "order_serive",
	})
}

func (h *Handler) GetMyCartHandler(c *gin.Context) {
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}
	data, err := h.services.Cart.CartList(user)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListCartResponse{Data: data})

}

func (h *Handler) AddProductHandler(c *gin.Context) {
	var input orderservice.CreateCartSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	data, err := h.services.Cart.CreateCart(user, input.Product)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, data)

}

func (h *Handler) UpdateCartHandler(c *gin.Context) {
	var input orderservice.UpdateCartSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	cart_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}

	access := h.services.Cart.UserCartPermission(user, cart_id)
	if !access {
		newErrorMessage(c, http.StatusForbidden, "no access to object")
		return
	}

	err = h.services.Cart.UpdateCart(cart_id, input.Amount)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) CleanCartHandler(c *gin.Context) {
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.services.Cart.CleanCart(user)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

func (h *Handler) RemoveCartPointHandler(c *gin.Context) {
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	cart_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	access := h.services.Cart.UserCartPermission(user, cart_id)
	if !access {
		newErrorMessage(c, http.StatusForbidden, "no access to object")
		return
	}

	err = h.services.Cart.RemoveCartPoint(cart_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
