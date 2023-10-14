package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ExposeHealthChecks включает ответчики для проверки работоспособности сервиса
func (tr *Transport) ExposeHealthChecks() {
	tr.Engine.GET("/ping", func(c *gin.Context) {
		tr.CounterService.Increment(c.Request.Context(), "ping_http", 1)
		c.AbortWithStatus(http.StatusNoContent)
	})
	tr.Engine.GET("/healthcheck", func(c *gin.Context) {
		ctx2, span := tr.SecretService.Tracer.Start(c.Request.Context(), "transport/http/healthcheck")
		defer span.End()
		tr.CounterService.Increment(ctx2, "healthcheck_http_called", 1)
		err := tr.SecretService.Ping(ctx2)
		if err != nil {
			tr.CounterService.Increment(ctx2, "healthcheck_http_failed", 1)
			log.Error().Err(err).
				Str("trace_id", span.SpanContext().TraceID().String()).
				Msgf("Ошибка при проверке сервиса: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		tr.CounterService.Increment(ctx2, "healthcheck_http_ok", 1)
		span.AddEvent("Сервис работает")
		log.Debug().
			Str("trace_id", span.SpanContext().TraceID().String()).
			Msgf("Сервис работает")
		c.String(http.StatusOK, "All systems online!")
	})

}
