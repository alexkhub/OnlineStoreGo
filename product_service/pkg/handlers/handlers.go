package handlers

import (
	"product_service/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct{
    services *service.Service
	auth service.JWTManager		
}

func NewHandler(services *service.Service, auth service.JWTManager) *Handler{
	return &Handler{services: services, auth: auth}
}

func (h *Handler) InitRouter() * gin.Engine{
	router:= gin.Default()
	router.GET("/", h.MainPage)
	router.MaxMultipartMemory = 15 << 20

	api := router.Group("/api")
	{
		product := api.Group("/product")
		{
			product.GET("/categories_list", h.ListCategoryHandler)
		}
		
		admin := api.Group("/admin", h.parseAuthHeader )
		{
			admin.POST("/create_category", h.CreateCategoryHanler)
			admin.POST("/create_product", h.CreateProductHandler)
			admin.POST("/upload_image/:id", h.AddImageHandler)
			
		}
		
	}
	
	
	return router
}