package handlers

import (
	
	"io"
	"net/http"
	productservice "product_service"
	"strconv"

	v "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategoryHanler(c *gin.Context){
    IsAdminPermission(c)
	
    var input productservice.CategorySerializer

    if err:= c.BindJSON(&input); err != nil{ 
        newErrorMessage(c, http.StatusBadRequest, err.Error())   
        return 
    }
    _, err := v.ValidateStruct(input)
    if err!= nil{
        newErrorMessage(c, http.StatusBadRequest, err.Error())
        return 
    }
    id, err := h.services.CreateCategory(input)

    if err!= nil{
        newErrorMessage(c, http.StatusInternalServerError, err.Error())
        return
    } 
    
    c.JSON(http.StatusOK, gin.H{
        "id": id,
    })
}

func (h *Handler) CreateProductHandler(c *gin.Context){
    IsAdminPermission(c)

    var input productservice.AdminCreateProductSerializer

     if err:= c.BindJSON(&input); err != nil{ 
        newErrorMessage(c, http.StatusBadRequest, err.Error())   
        return 
    }
    _, err := v.ValidateStruct(input)
    if err!= nil{
        newErrorMessage(c, http.StatusBadRequest, err.Error())
        return 
    }

    id, err := h.services.CreateProduct(input)

    if err != nil{
        newErrorMessage(c, http.StatusInternalServerError, err.Error())
        return 
    }

     c.JSON(http.StatusOK, gin.H{
        "id": id,
    })
}

func (h *Handler) AddImageHandler(c *gin.Context){
    IsAdminPermission(c)

    product_id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		newErrorMessage(c, http.StatusInternalServerError, err.Error())    
        return 
	}

    form, err := c.MultipartForm()
    if err != nil{
        newErrorMessage(c, http.StatusInternalServerError, err.Error())
        return  
    }
    files := form.File["files"]

    if files == nil{
        newErrorMessage(c, http.StatusBadRequest, "No files")

    }
    data := make(map[string]productservice.FileUploadSerializer)

    for _, file := range files{
        f, err := file.Open()
        if err != nil{
            newErrorMessage(c, http.StatusInternalServerError, err.Error())
            return  
        }
        defer f.Close()

        fileBytes, err := io.ReadAll(f)

        if err != nil{
            newErrorMessage(c, http.StatusInternalServerError, err.Error())
            return 
        }
        data[file.Filename] = productservice.FileUploadSerializer{FileName: file.Filename, Size: file.Size, Data: fileBytes}    
    }
    result, err := h.services.Admin.AddImage(product_id, data)

        if err != nil{
            newErrorMessage(c, http.StatusNotFound, err.Error())
            return 
        }
        
        c.JSON(http.StatusOK, result)  


}