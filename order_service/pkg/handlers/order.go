package handlers

import (
	"context"

	"net/http"
	orderservice "order_service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PaymentMethodeListHandler(c *gin.Context) {
	data, err := h.services.Order.PaymentMethodeList()
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListPaymentMethodeResponse{Data: data})
}

func (h *Handler) CreateOrderHandler(c *gin.Context) {
	var input orderservice.CreateOrderSerializer

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	input.User = user
	id, err := h.services.Order.CreateOrder(input)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})

}

func (h *Handler) DetailOrderHandler(c *gin.Context){
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan orderservice.MyError, 2)
	orderDataCh := make(chan orderservice.EmployeeOrderDataSerializer, 1)

	user, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}

	wg.Add(2)
	go func(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, orderId int64, userId int64){
		defer wg.Done()
		permission := h.services.Order.CheckOrderPermission(ctx, orderservice.OrderPermission{OrderId: int64(orderId), UserId: int64(user)})
		if permission.Error != nil{
			errCh <- permission
			ctxCancel()
			return
		}

	}(ctx, cancel, &wg, int64(orderId), int64(user))
	
	go func(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, orderId int64){
		defer wg.Done()
		data, err := h.services.Order.OrderDetail(ctx, orderId)
		if err != nil{
			errCh <- orderservice.MyError{
				Error: err,
				Code:  http.StatusInternalServerError,
			}
			ctxCancel()
			return
		}
		orderDataCh <- data
	}(ctx, cancel, &wg, int64(orderId))

	wg.Wait()
	close(errCh)
	close(orderDataCh)
	

	for opErr := range errCh {
		if opErr.Error != nil{
			newErrorMessage(c, opErr.Code, opErr.Error.Error())
			return
		}
	}
	
	data := <- orderDataCh

	c.JSON(http.StatusOK, data)

}


func (h *Handler) UserOrdersHandler(c *gin.Context) {
	user, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	data, err := h.services.Order.UserOrders(user)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, &getUserOrderListResponse{
		Data: data,
	})

}

func (h *Handler) OrderStatisticHandler(c *gin.Context) {
	user, err := GetUserId(c)
	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	data, err := h.services.Order.OrdersStatistic(user)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}