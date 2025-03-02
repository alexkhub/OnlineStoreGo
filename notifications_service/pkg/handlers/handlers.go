package handlers

import (
	"notifications_service/pkg/service"
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
	router.GET("/email",h.VerifyEmail)
	router.GET("/confirm/:uuid", h.AccountConfirm)

	return router
}