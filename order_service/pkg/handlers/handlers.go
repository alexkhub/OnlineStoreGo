package handlers

import (
	"order_service/pkg/service"

	"github.com/gin-gonic/gin"
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
		cart := api.Group("/cart", h.parseAuthHeader)
		{
			cart.GET("/my_cart", h.GetMyCartHandler)
			cart.POST("/add_product", h.AddProductHandler)
			cart.PATCH("/update_my_cart/:id", h.UpdateCartHandler)
			cart.DELETE("/clean_cart", h.CleanCartHandler)
			cart.DELETE("/delete_cart_point/:id", h.RemoveCartPointHandler)
		}
		order := api.Group("/order", h.parseAuthHeader)
		{
			order.GET("/payment_methode", h.PaymentMethodeListHandler)
			order.POST("/create_order", h.CreateOrderHandler)
		}
			
	}
	

	return router
}
