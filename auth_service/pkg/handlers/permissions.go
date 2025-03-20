package handlers

import (
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
)


func IsAuthenticatedPermission(c *gin.Context) (error){
	_, ok := c.Get("role")
	if !ok {
		return errors.New("role not found")
	}
	return nil
}

func IsStaffPermission(c *gin.Context) (error){
	
	role, ok := c.Get("role")
	if !ok {
		return errors.New("role not found")
	}
	roles := []string{"2", "4"}
	found := slices.Contains(roles, role.(string))
	if !found{
		return errors.New("authentication credentials were not provided")
	}
	return nil
}

func IsAdminPermission(c *gin.Context) (error){
	
	role, ok := c.Get("role")
	if !ok {
		return errors.New("role not found")
	}
	if role.(string) != "4"{
		return errors.New("authentication credentials were not provided")
	}
	return nil
}