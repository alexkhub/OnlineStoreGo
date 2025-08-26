package handlers

import (
	"context"
	"fmt"
	"net/http"
	orderservice "order_service"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ConfirmOrderStep1Handler(c *gin.Context) {
	IsStaffPermission(c)

	var wg sync.WaitGroup
	var input orderservice.ConfirmOrderStep1Serializer

	errCh := make(chan error, 2)
	resultCh := make(chan orderservice.EmployeeOrderDataSerializer, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	order_id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	input.Order = int64(order_id)

	wg.Add(2)
	go func(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, input orderservice.ConfirmOrderStep1Serializer) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			ctxCancel()
			return
		default:
		}

		err := h.services.Employee.ConfirmOrderStep1(ctx, input)
		if err != nil {
			errCh <- err
			ctxCancel()
			return
		}
	}(ctx, cancel, &wg, input)

	go func(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, orderId int64) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			ctxCancel()
			return
		default:
		}
		orderDetail, err := h.services.Order.OrderDetail(ctx, input.Order)
		if err != nil {
			errCh <- err
			ctxCancel()
			return
		}

		resultCh <- orderDetail

	}(ctx, cancel, &wg, input.Order)

	wg.Wait()
	close(errCh)
	close(resultCh)
	var errs []string
	for opErr := range errCh {
		errs = append(errs, opErr.Error())
	}

	if len(errs) != 0 {
		newErrorMessage(c, http.StatusInternalServerError, fmt.Sprintf("errors: %s", strings.Join(errs, ", ")))
		return
	}
	data := <-resultCh

	c.JSON(http.StatusOK, data)

}

func (h *Handler) ConfirmOrderStep2Handler(c *gin.Context) {
	IsStaffPermission(c)

	var input orderservice.UpdateListOrderPointSerializer
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	input.OrderId = int64(orderId)

	err = h.services.Employee.ConfirmOrderStep2(input)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)

}

func (h *Handler)ConfirmOrderStep3Handler(c *gin.Context) {
	IsStaffPermission(c)

	var input orderservice.ConfirmOrderStep3Serializer
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := GetUserId(c)

	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}

	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		newErrorMessage(c, http.StatusNotFound, err.Error())
		return
	}
	input.OrderId = int64(orderId)
	input.Employee = int64(user)
	err = h.services.Employee.ConfirmOrderStep3(context.Background(), input)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)




}