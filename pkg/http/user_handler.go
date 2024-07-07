package http

import (
	"fliqt/pkg/http/api"
	"fliqt/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *httpHandler) UserList(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))

	users, err := h.userSvc.List(ctx, page, size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (h *httpHandler) UserGet(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("userID")

	user, err := h.userSvc.FindByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *httpHandler) UserCreate(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.UserCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	user, err := h.userSvc.Create(ctx, model.User{
		ID:         req.ID,
		Name:       req.Name,
		Password:   req.Password,
		Role:       req.Role,
		Department: req.Department,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *httpHandler) UserDelete(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("userID")

	err := h.userSvc.DeleteByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func (h *httpHandler) UserUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("userID")

	var req model.User
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	user, err := h.userSvc.Update(ctx, userID, model.User{
		Name:       req.Name,
		Role:       req.Role,
		Department: req.Department,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
