package handlers

import (
	"github.com/gin-gonic/gin"
	"product_service/pkg/service"
)

type Handler struct {
	services *service.Service
	auth     service.JWTManager
}

func NewHandler(services *service.Service, auth service.JWTManager) *Handler {
	return &Handler{services: services, auth: auth}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", h.MainPage)
	router.MaxMultipartMemory = 15 << 20

	api := router.Group("/api")
	{
		product := api.Group("/product")
		{
			product.GET("/categories_list", h.ListCategoryHandler)
			product.GET("/product_list", h.ListProductHandler)
			product.GET("/product_detail/:id", h.ProductDetailHandler)
		}

		admin := api.Group("/admin", h.parseAuthHeader)
		{
			admin.POST("/create_category", h.CreateCategoryHanler)
			admin.POST("/create_product", h.CreateProductHandler)
			admin.POST("/upload_image/:id", h.AddImageHandler)

			admin.GET("/product_detail/:id", h.AdminProductDetailHandler)
			admin.PATCH("/product_detail/:id", h.UpdateProductHandler)
			admin.DELETE("/product_detail/:id", h.ProductDeleteHandler)
			admin.DELETE("/delete_image/:product_id/:name", h.RemoveImageHandler)

			admin.DELETE("/comment_remove/:id", h.AdminRemoveCommentHandler)

		}
		comment := api.Group("/comment")
		{
			comment.POST("/create_comment/:product_id", h.parseAuthHeader, h.CreateCommentHandler)
			comment.GET("/comment_list/:product_id", h.CommentListHandler)
			comment.DELETE("/comment_remove/:id", h.parseAuthHeader, h.CommentRemoveHandler)
		}

	}

	return router
}
