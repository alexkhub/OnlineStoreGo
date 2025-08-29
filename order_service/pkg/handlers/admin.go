package handlers

import (
	
	"net/http"
	"github.com/gin-gonic/gin"
)


func (h *Handler) AdminOrderListHandler(c *gin.Context){
	IsAdminPermission(c)

	filter, err := GenerateAdminFilter(c)
	if err != nil{
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	
	data, err := h.services.Admin.OrderList(filter)
	if err != nil{
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAdminOrderListResponse{Data: data})

}


func (h *Handler) AdminOrdersStatisticHandler(c *gin.Context){
	IsAdminPermission(c)

	filter, err := GenerateAdminFilter(c)
	if err != nil{
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	
	data, err := h.services.Admin.OrdersStatistic(filter)
	if err != nil{
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(data) == 1{
		c.JSON(http.StatusOK, data[0])
		return
	}

	c.JSON(http.StatusOK, getAdminOrderStatisticResponse{Data: data})

}