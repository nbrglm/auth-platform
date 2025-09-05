// Package tracing consists traces initialization and configures sending traces to Grafana Alloy.
//
// This package uses the opentelemetry sdks to configure otlp exporters over HTTP protocol to Grafana Alloy.
package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/nexeres/config"
	"github.com/nbrglm/nexeres/internal/logging"
	"github.com/nbrglm/nexeres/opts"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdouttrace "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	trace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// The global tracer instance.
var Tracer trace.Tracer

// The global tracer provider instance.
var Provider *sdktrace.TracerProvider

// Denotes an error during the initialization & configuration of the tracing system.
type TracingConfigurationError struct {
	Message string
}

func (e *TracingConfigurationError) Error() string {
	return e.Message
}

func InitTracer(ctx context.Context) (err error) {
	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(opts.Name), semconv.ServiceVersion(opts.Version), semconv.ServiceInstanceID(config.Server.InstanceID), semconv.DeploymentEnvironment(config.Environment()))

	var exporter sdktrace.SpanExporter

	switch config.Observability.OtelExporterProtocol {
	case "http/protobuf":
		options := []otlptracehttp.Option{
			otlptracehttp.WithEndpointURL(config.Observability.OtelExporterEndpoint),
		}
		if opts.Debug {
			options = append(options, otlptracehttp.WithInsecure()) // Use insecure connection in debug mode.
		}
		exporter, err = otlptracehttp.New(ctx, options...)
		if err != nil {
			return
		}
	case "grpc":
		options := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(config.Observability.OtelExporterEndpoint),
		}
		if opts.Debug {
			options = append(options, otlptracegrpc.WithInsecure()) // Use insecure connection in debug mode.
		}
		exporter, err = otlptracegrpc.New(ctx, options...)
		if err != nil {
			return
		}
	default:
		logging.Logger.Error("Unknown OTEL exporter protocol", zap.String("protocol", config.Observability.OtelExporterProtocol))
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return
		}
	}

	Provider = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithSampler(sdktrace.AlwaysSample()), sdktrace.WithResource(res))
	otel.SetTracerProvider(Provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Create the global tracer instance
	Tracer = otel.Tracer(opts.FullName)
	return
}

func ShutdownTracer(ctx context.Context) error {
	if Provider != nil {
		return Provider.Shutdown(ctx)
	}
	return nil
}

func AddTracingMiddleware(engine *gin.Engine) {
	engine.Use(otelgin.Middleware(opts.Name))
}
