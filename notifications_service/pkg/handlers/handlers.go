package handlers

import (
	"github.com/gin-gonic/gin"
	"notifications_service/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", h.MainPage)
	// router.GET("/email",h.VerifyEmail)
	router.GET("/confirm/:uuid", h.AccountConfirm)
	router.GET("/order_qr/:uuid", h.CheckQRHandler)

	return router
}
