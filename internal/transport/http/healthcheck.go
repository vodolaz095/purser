package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) ExposeHealthChecks() {
	tr.Engine.GET("/ping", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
	tr.Engine.GET("/healthcheck", func(c *gin.Context) {
		err := tr.Service.Ping(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "All systems online!")
	})

}
