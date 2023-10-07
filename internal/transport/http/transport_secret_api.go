package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/internal/transport/http/middlewares"
	"github.com/vodolaz095/purser/model"
)

type createSecretRequest struct {
	Body string            `json:"body" binding:"required"`
	Meta map[string]string `json:"meta"`
}

func makeLogger(c *gin.Context) zerolog.Logger {
	subj, found := c.Get("subject")
	if found {
		return log.With().
			Str("method", c.Request.Method).
			Str("endpoint", c.Request.RequestURI).
			Str("remote_addr", c.RemoteIP()).
			Str("subject", subj.(string)).
			Str("user_agent", c.GetHeader("User-Agent")).
			Logger()
	}
	return log.With().
		Str("method", c.Request.Method).
		Str("endpoint", c.Request.RequestURI).
		Str("remote_addr", c.RemoteIP()).
		Str("user_agent", c.GetHeader("User-Agent")).
		Logger()
}

func (tr *Transport) ExposeSecretAPI() {
	rest := tr.Engine.Group("/api/v1/secret")
	rest.Use(middlewares.CheckJWT())

	rest.GET("/", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotImplemented)
	})
	rest.PUT("/:id", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotImplemented)
	})
	rest.PATCH("/:id", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotImplemented)
	})
	rest.POST("/:id", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotImplemented)
	})
	rest.GET("/:id", func(c *gin.Context) {
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/GetSecretByID")
		defer span.End()
		logger := makeLogger(c)
		id := c.Param("id")
		secret, err := tr.Service.FindByID(ctx2, id)
		if err != nil {
			if errors.Is(err, model.SecretNotFoundError) {
				logger.Info().
					Str("trace_id", span.SpanContext().TraceID().String()).
					Str("secret_id", id).
					Msgf("Секрет %s не найден", id)
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			logger.Error().Err(err).
				Str("trace_id", span.SpanContext().TraceID().String()).
				Str("secret_id", id).
				Msgf("Ошибка при поиске секреа %s : %s", id, err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Info().
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("secret_id", id).
			Msgf("Секрет %s получен", id)
		c.JSON(http.StatusOK, secret)
	})
	rest.DELETE("/:id", func(c *gin.Context) {
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/DeleteSecretByID")
		defer span.End()
		logger := makeLogger(c)
		id := c.Param("id")
		err := tr.Service.DeleteByID(ctx2, c.Param("id"))
		if err != nil {
			if errors.Is(err, model.SecretNotFoundError) {
				logger.Info().
					Str("trace_id", span.SpanContext().TraceID().String()).
					Str("secret_id", id).
					Msgf("Секрет %s не найден", id)
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			logger.Error().Err(err).
				Str("trace_id", span.SpanContext().TraceID().String()).
				Str("secret_id", id).
				Msgf("Ошибка при поиске секрета %s : %s", id, err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Info().
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("secret_id", id).
			Msgf("Секрет %s удалён", id)
		c.AbortWithStatus(http.StatusNoContent)
	})
	rest.POST("/", func(c *gin.Context) {
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/CreateSecret")
		defer span.End()
		logger := makeLogger(c)
		subject, found := c.Get("subject")
		if !found {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var bdy createSecretRequest
		if err := c.ShouldBindJSON(&bdy); err != nil {
			logger.Info().Err(err).
				Str("trace_id", span.SpanContext().TraceID().String()).
				Msgf("Ошибка при валидации секрета: %s", err)
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
		bdy.Meta["Subject"] = subject.(string)
		secret, err := tr.Service.Create(ctx2, bdy.Body, bdy.Meta)
		if err != nil {
			logger.Error().Err(err).
				Str("trace_id", span.SpanContext().TraceID().String()).
				Msgf("Ошибка при создании секрета: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Info().Err(err).
			Str("secret_id", secret.ID).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Msgf("Пользователь %s создал секрет %s", subject.(string), secret.ID)
		c.Header("Location", "/api/v1/secret/"+secret.ID)
		c.AbortWithStatus(http.StatusCreated)
	})
}
