package handlers


import (
	"auth_service/pkg/service"
	"github.com/gin-gonic/gin"
)


type Handler struct{
    services *service.Service
		
}

func NewHandler(services *service.Service) *Handler{
	return &Handler{services: services}
}


func (h *Handler) InitRouter() * gin.Engine{
	router:= gin.Default()
	router.GET("/", h.MainPage)
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/registration", h.RegistrationHandler)		
			auth.POST("/login", h.LoginHandler)
		}
		
	}
	return router
}