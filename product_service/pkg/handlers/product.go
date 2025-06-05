package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)



func (h *Handler) MainPage(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{
        "message": "product_serive",
    })
}

func (h *Handler) ListCategoryHandler(c *gin.Context){
    data, err := h.services.Product.CatregoList()
    if err != nil{
        newErrorMessage(c, http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, getListCategoryResponse{Data: data})
}