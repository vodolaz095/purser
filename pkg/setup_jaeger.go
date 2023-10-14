package pkg

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// SetupJaeger внезапно, настраивает систему отправки телеметрии
func SetupJaeger(hostname, version, environment, jaegerHost, jaegerPort string) error {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(jaegerHost),
		jaeger.WithAgentPort(jaegerPort),
	))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// sample 10% of data to save bandwidth
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(0.1)),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("purser"),
			semconv.HostNameKey.String(hostname),
			semconv.DeploymentEnvironmentKey.String(environment),
			semconv.ServiceVersion(version),
		)),
	)
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	return nil
}
