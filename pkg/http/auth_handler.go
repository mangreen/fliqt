package http

import (
	"encoding/json"
	"fliqt/pkg/common"
	"fliqt/pkg/http/api"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *httpHandler) AuthLogin(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	usr, err := h.authSvc.Validate(ctx, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	jsonBytes, err := json.Marshal(usr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	ss := sessions.Default(c)

	ss.Set(common.FLIQT_CONST, string(jsonBytes))

	ss.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})
	// sessions.Session.Save() is for setting cookie
	if err := ss.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	// /* another way to set cookie */
	// ssID := ss.ID()
	// c.SetCookie(common.FLIQT_CONST, "session_"+ssID, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"data": usr,
	})
}

func (h *httpHandler) AuthLogout(c *gin.Context) {
	ss := sessions.Default(c)

	ss.Delete(common.FLIQT_CONST)

	// set maxAge -1 to delete cookie
	ss.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	if err := ss.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	// /* another way to get & delete cookie */
	// _, err := c.Cookie(common.FLIQT_CONST)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"err": err.Error(),
	// 	})
	// 	return
	// }
	// c.SetCookie(common.FLIQT_CONST, "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"msg": "logout success",
	})
}
