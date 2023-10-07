package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (tr *Transport) ExposeHealthChecks() {
	tr.Engine.GET("/ping", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
	tr.Engine.GET("/healthcheck", func(c *gin.Context) {
		ctx2, span := tr.Service.Tracer.Start(c.Request.Context(), "transport/http/healthcheck")
		defer span.End()
		err := tr.Service.Ping(ctx2)
		if err != nil {
			log.Error().Err(err).
				Str("trace_id", span.SpanContext().TraceID().String()).
				Msgf("Ошибка при проверке сервиса: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		span.AddEvent("Сервис работает")
		log.Debug().
			Str("trace_id", span.SpanContext().TraceID().String()).
			Msgf("Сервис работает")
		c.String(http.StatusOK, "All systems online!")
	})

}
