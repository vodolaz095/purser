package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/pkg"
)

func CheckJWT(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.String(http.StatusUnauthorized, "Authorization header required")
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.String(http.StatusUnauthorized, "Authorization header with bearer strategy required")
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		subject, err := pkg.ValidateJwtAndExtractSubject(token, config.JwtSecret)
		if err != nil {
			c.String(http.StatusBadRequest, "Error parsing token: %s", err)
			return
		}
		c.Set("subject", subject)
		c.Next()
	})
}
