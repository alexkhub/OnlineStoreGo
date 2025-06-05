package handlers

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)


func IsAuthenticatedPermission(c *gin.Context) {
	_, ok := c.Get("role")
	if !ok {
		newErrorMessage(c, http.StatusForbidden, errors.New("role not found").Error())
		return
	}
}

func IsStaffPermission(c *gin.Context) {
	
	role, ok := c.Get("role")
	if !ok {
		newErrorMessage(c, http.StatusForbidden, errors.New("role not found").Error())
		return
	}
	roles := []string{"2", "4"}
	found := slices.Contains(roles, role.(string))
	if !found{
		newErrorMessage(c, http.StatusForbidden, errors.New("authentication credentials were not provided").Error())    
        return 
		
	}
	
}

func IsAdminPermission(c *gin.Context){
	
	role, ok := c.Get("role")
	if !ok {
		newErrorMessage(c, http.StatusForbidden, errors.New("role not found").Error())
		return
	}
	if role.(string) != "4"{
		newErrorMessage(c, http.StatusForbidden, errors.New("authentication credentials were not provided").Error())    
        return 
	}
}