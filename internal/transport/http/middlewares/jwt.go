package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/pkg"
)

func CheckJWT() func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.String(http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.String(http.StatusUnauthorized, "Authorization header with bearer strategy required")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		subject, err := pkg.ValidateJwtAndExtractSubject(token, config.JwtSecret)
		if err != nil {
			c.String(http.StatusBadRequest, "Error parsing token: %s", err)
			c.Abort()
			return
		}
		// также тут можно вызывать некое хранилище отозванных токенов, чтобы проверить,
		// что этот ещё работает, а также можно проверить роли и разрешения пользователя
		c.Set("subject", subject)
		c.Next()
	}
}
