package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckQRHandler(c *gin.Context) {
	err := h.services.OrderConfirmStep1(c.Param("uuid"))
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "account_confirm",
	})
}
