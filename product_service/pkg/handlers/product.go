package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MainPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "product_serive",
	})
}

func (h *Handler) ListCategoryHandler(c *gin.Context) {
	data, err := h.services.Product.CategoryList()
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListCategoryResponse{Data: data})
}

func (h *Handler) ListProductHandler(c *gin.Context) {
	data, err := h.services.Product.ProductList()
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListProductResponse{Data: data})
}

func (h *Handler) ProductDetailHandler(c *gin.Context) {
	product_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	chech_product := h.services.Product.CheckProduct(product_id)
	if !chech_product {
		newErrorMessage(c, http.StatusNotFound, "object not found")
		return
	}

	product, err := h.services.ProductDetail(product_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, product)
}
