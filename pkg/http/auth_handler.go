package http

import (
	"fliqt/pkg/http/api"
	"net/http"

	// "fliqt/pkg/delivery/http/model"

	"github.com/gin-gonic/gin"
)

func (h *httpHandler) AuthLogin(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})

		return
	}

	res, err := h.authSvc.Validate(ctx, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
