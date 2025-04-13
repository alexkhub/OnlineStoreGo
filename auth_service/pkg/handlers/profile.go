package handlers

import (
	authservice "auth_service"
	"io"
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

func (h *Handler) ProfileUploadFileHandler (c *gin.Context){
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
	file, err := c.FormFile("file")

	if err != nil{
		newErrorMessage(c, http.StatusBadRequest, err.Error())
        return
	}

	f, _ := file.Open()
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil{
		newErrorMessage(c, http.StatusBadRequest, err.Error())
        return
	}



	err = h.services.Profile.UpdateProfileImage(user, authservice.FileUploadSerializer{FileName: file.Filename, Size: file.Size, Data: fileBytes})

	if err != nil{
		newErrorMessage(c, http.StatusBadRequest, err.Error())
        return
	}

	c.JSON(http.StatusOK, nil)
}