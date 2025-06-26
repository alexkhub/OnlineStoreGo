package handlers

import (
	"io"
	"log"
	"net/http"
	productservice "product_service"
	"strconv"

	v "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategoryHanler(c *gin.Context) {
	IsAdminPermission(c)

	var input productservice.CategorySerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateCategory(input)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) CreateProductHandler(c *gin.Context) {
	IsAdminPermission(c)

	var input productservice.AdminCreateProductSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateProduct(input)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) AddImageHandler(c *gin.Context) {
	IsAdminPermission(c)

	product_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	files := form.File["files"]

	if files == nil {
		newErrorMessage(c, http.StatusBadRequest, "No files")
		return

	}
	data := make(map[string]productservice.FileUploadSerializer)

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			newErrorMessage(c, http.StatusInternalServerError, err.Error())
			return
		}
		defer f.Close()
		fileBytes, err := io.ReadAll(f)
		if err != nil {
			newErrorMessage(c, http.StatusInternalServerError, err.Error())
			return
		}
		data[file.Filename] = productservice.FileUploadSerializer{FileName: file.Filename, Size: file.Size, Data: fileBytes}
	}
	result, err := h.services.Admin.AddImage(product_id, data)

	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

}

func (h *Handler) ProductDeleteHandler(c *gin.Context) {
	IsAdminPermission(c)
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

	err = h.services.ProductDelete(product_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) AdminProductDetailHandler(c *gin.Context) {
	IsAdminPermission(c)
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

	data, err := h.services.AdminProductDetail(product_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) RemoveImageHandler(c *gin.Context) {
	IsAdminPermission(c)
	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	image := c.Param("name")
	err = h.services.RemoveImage(product_id, image)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateProductHandler(c *gin.Context) {
	IsAdminPermission(c)
	var input productservice.AdminUpdateProductSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := v.ValidateStruct(input)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}

	log.Println(input.Category)

	err = h.services.UpdateProduct(id, input)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := h.services.AdminProductDetail(id)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)

}
