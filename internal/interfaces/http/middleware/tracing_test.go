package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestTracingMiddleware_CreatesServerSpan(t *testing.T) {
	tracerProvider := sdktrace.NewTracerProvider()

	t.Cleanup(func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			t.Errorf("shutdown tracer provider: %v", err)
		}
	})

	handler := TracingMiddleware(tracerProvider)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanContext := trace.SpanContextFromContext(r.Context())

			if !spanContext.IsValid() {
				t.Fatal("expected valid span context")
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/photos",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			rec.Code,
		)
	}
}

func TestTracingMiddleware_ExportsServerSpan(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(recorder),
	)

	t.Cleanup(func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			t.Errorf("shutdown tracer provider: %v", err)
		}
	})

	handler := TracingMiddleware(tracerProvider)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/photos",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	spans := recorder.Ended()

	if len(spans) != 1 {
		t.Fatalf(
			"expected 1 ended span, got %d",
			len(spans),
		)
	}

	span := spans[0]

	if !span.SpanContext().IsValid() {
		t.Fatal("expected exported span to have a valid span context")
	}

	if span.SpanKind() != trace.SpanKindServer {
		t.Fatalf(
			"expected server span kind, got %v",
			span.SpanKind(),
		)
	}
}

func TestTracingMiddleware_ExtractsIncomingTraceContext(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(recorder),
	)

	t.Cleanup(func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			t.Errorf("shutdown tracer provider: %v", err)
		}
	})

	previousPropagator := otel.GetTextMapPropagator()

	t.Cleanup(func() {
		otel.SetTextMapPropagator(previousPropagator)
	})

	propagator := propagation.TraceContext{}

	otel.SetTextMapPropagator(propagator)

	traceID, err := trace.TraceIDFromHex(
		"4bf92f3577b34da6a3ce929d0e0e4736",
	)
	if err != nil {
		t.Fatalf("create trace ID: %v", err)
	}

	spanID, err := trace.SpanIDFromHex(
		"00f067aa0ba902b7",
	)
	if err != nil {
		t.Fatalf("create span ID: %v", err)
	}

	parentSpanContext := trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		},
	)

	parentContext := trace.ContextWithRemoteSpanContext(
		context.Background(),
		parentSpanContext,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/photos",
		nil,
	)

	propagator.Inject(
		parentContext,
		propagation.HeaderCarrier(req.Header),
	)

	handler := TracingMiddleware(tracerProvider)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanContext := trace.SpanContextFromContext(r.Context())

			if !spanContext.IsValid() {
				t.Fatal("expected valid span context")
			}

			if spanContext.TraceID() != traceID {
				t.Fatalf(
					"expected trace ID %s, got %s",
					traceID,
					spanContext.TraceID(),
				)
			}

			if spanContext.SpanID() == spanID {
				t.Fatal("expected server span to have a different span ID")
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			rec.Code,
		)
	}

	spans := recorder.Ended()

	if len(spans) != 1 {
		t.Fatalf(
			"expected 1 ended span, got %d",
			len(spans),
		)
	}

	if spans[0].Parent().SpanID() != spanID {
		t.Fatalf(
			"expected parent span ID %s, got %s",
			spanID,
			spans[0].Parent().SpanID(),
		)
	}
}

func TestSpanNameFormatter(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		pattern string
		want    string
	}{
		{
			name:    "route pattern",
			method:  http.MethodGet,
			pattern: "/api/v1/photos/{id}",
			want:    "GET /api/v1/photos/{id}",
		},
		{
			name:   "method fallback",
			method: http.MethodPost,
			want:   "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				"/api/v1/photos/123",
				nil,
			)

			req.Pattern = tt.pattern

			got := spanNameFormatter("", req)

			if got != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}
