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
	router.MaxMultipartMemory = 15 << 20
	
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/registration", h.RegistrationHandler)		
			auth.POST("/login", h.LoginHandler)
			auth.POST("/refresh", h.RefreshJWTHandler)
			auth.POST("/logout", h.LogoutHandler )
			auth.GET("/close_all_sessions", h.parseAuthHeader,  h.CloseAllSessionsHandler)
		}
		profile := api.Group("/profile", h.parseAuthHeader)
		{
			
			profile.GET("/", h.ProfileHandler)
			profile.POST("/upload_img", h.ProfileUploadFileHandler)
			profile.PATCH("/", h.ProfileUpdateHandler)
			profile.DELETE("/", h.ProfileDeleteHandler)
		}
		admin := api.Group("/admin", h.parseAuthHeader )
		{
			admin.GET("/user_list", h.UserListHandler)
			admin.GET("/role_list", h.RoleListHandker)
			admin.GET("user/:id", h.UserDetailHandler)
			admin.PATCH("user/:id", h.UserUpdateHandler)
			admin.DELETE("user/:id", h.UserDeleteHandler)
			admin.POST("unblock/:id", h.UserUnblockHandler)
			
		}
		
	}
	return router
}