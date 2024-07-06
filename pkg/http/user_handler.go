package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *httpHandler) UserGet(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("userID")

	usr, err := h.userSvc.FindByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": usr,
	})
}
