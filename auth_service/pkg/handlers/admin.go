package handlers

import (
	authservice "auth_service"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UserListHandler(c *gin.Context) {
	IsStaffPermission(c)
	filter := c.Request.URL.Query()
	user_list, err := h.services.UserList(filter)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListUserResponse{
		Data: user_list,
	})

}

func (h *Handler) RoleListHandker(c *gin.Context) {
	IsStaffPermission(c)

	user_role, err := h.services.RoleList()

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListRoleResponse{
		Data: user_role,
	})
}

func (h *Handler) UserDetailHandler(c *gin.Context) {
	IsStaffPermission(c)

	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}

	data, err := h.services.Profile.UserProfile(user_id)

	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)

}

func (h *Handler) UserUpdateHandler(c *gin.Context) {
	IsAdminPermission(c)

	var input authservice.ProfileSerializer
	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusForbidden, err.Error())
		return
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	
	err = h.services.ProfileUpdate(user_id, input)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile update"})
}

func (h *Handler) UserDeleteHandler(c *gin.Context) {
	IsAdminPermission(c)

	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.Admin.UserBlock(user_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user block"})

}

func (h *Handler) UserUnblockHandler(c *gin.Context) {
	IsAdminPermission(c)

	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.Admin.UserUnblock(user_id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user unblock"})
}
