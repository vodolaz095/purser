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
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/GetSecretByID")
		defer span.End()
		secret, err := tr.Service.FindByID(ctx2, c.Param("id"))
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
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/DeleteSecretByID")
		defer span.End()
		err := tr.Service.DeleteByID(ctx2, c.Param("id"))
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
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/CreateSecret")
		defer span.End()
		var found bool
		var bdy createSecretRequest
		if err := c.ShouldBindJSON(&bdy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if bdy.Meta == nil {
			bdy.Meta = make(map[string]string, 0)
		}
		_, found = bdy.Meta["body"]
		if found {
			delete(bdy.Meta, "body")
		}
		bdy.Meta["User-Agent"] = c.Request.Header.Get("User-Agent")
		secret, err := tr.Service.Create(ctx2, bdy.Body, bdy.Meta)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Header("Location", "/api/v1/secret/"+secret.ID)
		c.AbortWithStatus(http.StatusCreated)
	})
}
