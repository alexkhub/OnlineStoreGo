package handlers

import (
	"errors"
	"net/http"
	orderservice "order_service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) parseAuthHeader(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorMessage(c, http.StatusForbidden, errors.New("empty auth header").Error())
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorMessage(c, http.StatusForbidden, errors.New("invalid auth header").Error())
		return
	}

	if len(headerParts[1]) == 0 {

		newErrorMessage(c, http.StatusForbidden, errors.New("token is empty").Error())
		return
	}
	res, err := h.auth.Parse(headerParts[1])

	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}
	c.Set("user_id", res.Id)
	c.Set("role", res.Role)

}

func GetUserId(c *gin.Context) (int, error) {

	id, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, err := strconv.Atoi(id.(string))
	if err != nil {
		return 0, errors.New("user id not convert")
	}
	return idInt, nil
}


func GenerateAdminFilter(c *gin.Context) (orderservice.OrderFilter, error){
	var filter orderservice.OrderFilter

	filterParam := c.Request.URL.Query()

	if !filterParam.Has("create_at_gte"){
		filter.CreateAtGTE = time.Now().Format(time.DateOnly)
	}else{
		CreateAtGTE, err := time.Parse(time.DateOnly, filterParam.Get("create_at_gte"))
		if err != nil{
			return orderservice.OrderFilter{}, err
		}
		filter.CreateAtGTE = CreateAtGTE.Format(time.DateOnly)
	}
	if filterParam.Has("create_at_lte"){
		CreateAtLTE, err := time.Parse(time.DateOnly, filterParam.Get("create_at_lte"))
		if err != nil{
			return orderservice.OrderFilter{}, err
		}
		filter.CreateAtGTE = CreateAtLTE.Format(time.DateOnly)
	}

	if filterParam.Has("min_price"){
		minPrice, err := strconv.Atoi(filterParam.Get("min_price"))
		if err != nil{
			return orderservice.OrderFilter{}, err
		}
		filter.MinPrice.SetValid(int64(minPrice))
	}
	if filterParam.Has("max_price"){
		maxPrice, err := strconv.Atoi(filterParam.Get("max_price"))
		if err != nil{
			return orderservice.OrderFilter{}, err
		}
		filter.MaxPrice.SetValid(int64(maxPrice))
	}
	if filterParam.Has("payment_method"){
		filter.PaymentMethode.SetValid(filterParam.Get("payment_method"))
	}
	if filterParam.Has("status"){
		filter.Status.SetValid(filterParam.Get("status"))
	}

	return filter, nil
}