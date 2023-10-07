package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/purser/model"
)

type createSecretRequest struct {
	Body string            `json:"body" binding:"required"`
	Meta map[string]string `json:"meta"`
}

func (tr *Transport) ExposeSecretAPI() {
	rest := tr.Engine.Group("/api/v1/secret")
	rest.GET("/", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotImplemented)
	})
	rest.GET("/:id", func(c *gin.Context) {
		secret, err := tr.Service.FindByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			if errors.Is(err, model.SecretNotFoundError) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, secret)
	})
	rest.DELETE("/:id", func(c *gin.Context) {
		err := tr.Service.DeleteByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			if errors.Is(err, model.SecretNotFoundError) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatus(http.StatusNoContent)
	})
	rest.POST("/", func(c *gin.Context) {
		var bdy createSecretRequest
		if err := c.ShouldBindJSON(&bdy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		secret, err := tr.Service.Create(c.Request.Context(), bdy.Body, bdy.Meta)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Header("Location", "/api/v1/secret/"+secret.ID)
		c.AbortWithStatus(http.StatusCreated)
	})
}
