package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// UseTracing adds open telemetry tracing
func UseTracing(router *gin.Engine) {
	router.Use(otelgin.Middleware("purser",
		otelgin.WithSpanNameFormatter(func(r *http.Request) string {
			return r.Method + " " + r.URL.Path
		})))
}
