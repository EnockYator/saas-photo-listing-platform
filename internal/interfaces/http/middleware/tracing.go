package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

const tracingOperation = "http.server"

// TracingMiddleware instruments inbound HTTP requests with OpenTelemetry.
//
// It:
//   - creates a server span for each traced request;
//   - extracts incoming trace context;
//   - propagates the active span through request.Context();
//   - records standard HTTP telemetry;
//   - uses the configured TracerProvider when one is supplied.
//
// OpenTelemetry owns trace and span context. This middleware intentionally
// does not copy trace IDs into the application's requestcontext package.
func TracingMiddleware(
	tracerProvider trace.TracerProvider,
) func(http.Handler) http.Handler {
	options := []otelhttp.Option{
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
	}

	if tracerProvider != nil {
		options = append(
			options,
			otelhttp.WithTracerProvider(tracerProvider),
		)
	}

	return otelhttp.NewMiddleware(
		tracingOperation,
		options...,
	)
}

// spanNameFormatter produces stable, low-cardinality span names.
//
// The route pattern is preferred when available. Falling back to the HTTP
// method avoids putting arbitrary URL paths into span names.
func spanNameFormatter(
	_ string,
	r *http.Request,
) string {
	if pattern := r.Pattern; pattern != "" {
		return r.Method + " " + pattern
	}

	return r.Method
}
