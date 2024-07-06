package http

import (
	"fliqt/pkg/http/middle"
	"fliqt/pkg/svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	httpHandler struct {
		userSvc svc.UserService
		authSvc svc.AuthService
	}
)

func NewHandler(r *gin.Engine, userSvc svc.UserService, authSvc svc.AuthService) {
	handler := &httpHandler{
		userSvc: userSvc,
		authSvc: authSvc,
	}

	r.GET("/", handler.Home)

	api := r.Group("/api")

	userApi := api.Group("/users")
	userApi.Use(middle.AuthMiddleware())
	userApi.GET("/:userID", handler.UserGet)

	authApi := api.Group("/auth")
	authApi.POST("/login", handler.AuthLogin)
}

func (h *httpHandler) Home(c *gin.Context) {
	c.String(http.StatusOK, "Welcome Gin Server")
}
