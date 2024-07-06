package middle

import (
	"encoding/json"
	"fliqt/pkg/common"
	"fliqt/pkg/model"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		ss := sessions.Default(c)

		jsonStr := ss.Get(common.FLIQT_CONST)
		if jsonStr == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "No seesion found",
			})
			return
		}

		usr := &model.User{}
		if err := json.Unmarshal([]byte(jsonStr.(string)), usr); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": err.Error(),
			})
			return
		}

		// Store User info into Context
		c.Set("user", usr)
		// After that, we can get User info from c.Get("user")
		c.Next()
	}
}
