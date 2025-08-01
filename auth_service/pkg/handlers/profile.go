package handlers

import (
	authservice "auth_service"
	"io"
	"net/http"

	v "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ProfileHandler(c *gin.Context) {
	IsAuthenticatedPermission(c)

	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	data, err := h.services.Profile.UserProfile(user)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) ProfileUploadFileHandler(c *gin.Context) {
	IsAuthenticatedPermission(c)
	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	file, err := c.FormFile("file")

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	f, _ := file.Open()
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Profile.UpdateProfileImage(user, authservice.FileUploadSerializer{FileName: file.Filename, Size: file.Size, Data: fileBytes})

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) ProfileUpdateHandler(c *gin.Context) {
	var input authservice.ProfileSerializer
	user, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err = v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.ProfileUpdate(user, input)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile update"})
}

func (h *Handler) ProfileDeleteHandler(c *gin.Context) {
	user, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.services.ProfileDelete(user)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile close"})

}
