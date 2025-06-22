package handlers

import (
	"net/http"

	"auth_service"

	v "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func (h *Handler) MainPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "auth_serive",
	})
}

func (h *Handler) RegistrationHandler(c *gin.Context) {
	var input authservice.AuthRegistrationSerializer

	if err := c.BindJSON(&input); err != nil {

		newErrorMessage(c, http.StatusBadRequest, err.Error())

		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())

		return
	}

	response, err := h.services.Registration(input)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var input authservice.LoginUser

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.LoginUser(input)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)

}

func (h *Handler) RefreshJWTHandler(c *gin.Context) {
	var input authservice.RefreshToken

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.RefreshJWTToken(input.Refresh)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)

}

func (h *Handler) LogoutHandler(c *gin.Context) {
	var input authservice.RefreshToken

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.DeleteRefreshJWTToken(input.Refresh)

	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})

}

func (h *Handler) CloseAllSessionsHandler(c *gin.Context) {
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.services.CloseAllSessions(user)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "close all sessions",
	})
}
