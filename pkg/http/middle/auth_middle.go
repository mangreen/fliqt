package middle

import (
	"errors"
	"fliqt/pkg/svc"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ParseToken(tokenString string) (*svc.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &svc.JwtClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(svc.FLIQT_CONST), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*svc.JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "Authorization is null in Header",
			})
			return
		}

		// Headers: {"Authorization": "Bearer <token>"}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "Format of Authorization is wrong",
			})
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "Invalid Token",
			})
			return
		}

		// Store User info into Context
		c.Set("user", claims.User)
		// After that, we can get User info from c.Get("user")
		c.Next()
	}
}
