package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	//  v "github.com/asaskevich/govalidator"
)


func (h *Handler) MainPage(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{
        "message": "notification_serive",
    })
}

func (h *Handler) AccountConfirm(c *gin.Context){
    err := h.services.AccountConfirm(c.Param("uuid"))
    if err!= nil{
        newErrorMessage(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "account_confirm",
    })
}