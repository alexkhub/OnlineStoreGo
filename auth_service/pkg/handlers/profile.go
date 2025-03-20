package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) ProfileHandler (c *gin.Context){
	err := IsAuthenticatedPermission(c)

	if err != nil{
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
        return
	}
	user, err := GetUserId(c)

	if err != nil{
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
        return
	}

	data, err := h.services.Profile.UserProfile(user)

	if err != nil{
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
        return
	}
	c.JSON(http.StatusOK, data)
}