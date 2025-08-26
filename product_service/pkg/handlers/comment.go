package handlers

import (
	"net/http"
	productservice "product_service"
	"strconv"


	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCommentHandler(c *gin.Context) {
	var input productservice.CreateCommentSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	user_id, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}
	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	chech_product := h.services.Product.CheckProduct(product_id)
	if !chech_product {
		newErrorMessage(c, http.StatusNotFound, "object not found")
		return
	}

	id, err := h.services.CreateComment(input, product_id, user_id)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})

}



func (h *Handler) CommentListHandler(c *gin.Context){
	
	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	chech_product := h.services.Product.CheckProduct(product_id)
	if !chech_product {
		newErrorMessage(c, http.StatusNotFound, "product not found")
		return
	}

	data, err := h.services.CommentList(product_id)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListCommentResponse{Data: data})

}


func (h *Handler) CommentRemoveHandler(c *gin.Context){
	comment, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	user_id, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}

	err = h.services.Comment.RemoveComment(comment, user_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}