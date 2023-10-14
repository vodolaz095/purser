package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/codes"
)

// https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example

// metricsToExpose задают метрики, которые экспортируются в Prometheus. Я их нашёл
// такой командой из кода `$ grep "CounterService.Increment" internal/transport/**/*.go`
var metricsToExpose = []string{
	"grpc_get_secret_called",
	"grpc_get_secret_not_found",
	"grpc_get_secret_error",
	"grpc_get_secret_success",
	"grpc_delete_secret_called",
	"grpc_delete_secret_not_found",
	"grpc_delete_secret_error",
	"grpc_delete_secret_success",
	"grpc_create_secret_called",
	"grpc_create_secret_error",
	"grpc_create_secret_success",
	"ping_http",
	"healthcheck_http_called",
	"healthcheck_http_failed",
	"healthcheck_http_ok",
	"http_get_secret_called",
	"http_get_secret_not_found",
	"http_get_secret_error",
	"http_get_secret_success",
	"http_delete_secret_called",
	"http_delete_secret_not_found",
	"http_delete_secret_error",
	"http_delete_secret_success",
	"http_create_secret_called",
	"http_create_secret_malformed",
	"http_create_secret_error",
	"http_create_secret_success",
}

func (tr *Transport) ExposeMetrics() {
	tr.Engine.GET("/metrics", func(c *gin.Context) {
		var err error
		ctx2, span := tr.SecretService.Tracer.Start(c.Request.Context(), "transport/http/metrics")
		defer span.End()
		c.Header("Content-Type", "text/plain; version=0.0.4")
		for i := range metricsToExpose {
			val, found := tr.CounterService.Get(ctx2, metricsToExpose[i])
			if found {
				_, err = fmt.Fprintf(c.Writer, "%s{hostname=\"%s\"} %.2f\n",
					metricsToExpose[i], tr.Hostname, float64(val),
				)
			} else {
				_, err = fmt.Fprintf(c.Writer, "%s{hostname=\"%s\"} 0\n",
					metricsToExpose[i], tr.Hostname)
			}
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				log.Error().Err(err).
					Msgf("ошибка отправки данных: %s", err)
				break
			}
		}
		span.AddEvent("Метрики отправлены")
		fmt.Fprintln(c.Writer)
		c.AbortWithStatus(http.StatusOK)
	})
}
