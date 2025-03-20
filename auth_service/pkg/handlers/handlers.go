package handlers


import (
	"auth_service/pkg/service"
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
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/registration", h.RegistrationHandler)		
			auth.POST("/login", h.LoginHandler)
			auth.POST("/refresh", h.RefreshJWTHandler)
		}
		profile := api.Group("/profile", h.parseAuthHeader)
		{
			profile.GET("/", h.ProfileHandler)
		}
		
	}
	return router
}