package metrics

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcMetrics struct {
	// Latency - Duration histogram for p50, p95, p99
	requestDuration metric.Float64Histogram

	// Traffic - Request counter
	requestsTotal metric.Int64Counter

	// Errors - Error counter by code
	errorsTotal metric.Int64Counter

	// Saturation - Active requests gauge
	activeRequests metric.Int64UpDownCounter
}

func NewGrpcMetrics(meter metric.Meter) (*GrpcMetrics, error) {
	gm := &GrpcMetrics{}

	var err error

	// Latency histogram with buckets optimized for p50, p95, p99
	gm.requestDuration, err = meter.Float64Histogram(
		"grpc.server.request_duration",
		metric.WithDescription("gRPC request duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			0.001, // 1ms
			0.005, // 5ms
			0.01,  // 10ms
			0.025, // 25ms
			0.05,  // 50ms
			0.1,   // 100ms
			0.25,  // 250ms
			0.5,   // 500ms
			1.0,   // 1s
			2.5,   // 2.5s
			5.0,   // 5s
			10.0,  // 10s
		),
	)
	if err != nil {
		return nil, err
	}

	// Traffic counter
	gm.requestsTotal, err = meter.Int64Counter(
		"grpc.server.requests_total",
		metric.WithDescription("Total number of gRPC requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	// Error counter
	gm.errorsTotal, err = meter.Int64Counter(
		"grpc.server.errors_total",
		metric.WithDescription("Total number of gRPC errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return nil, err
	}

	// Active requests (saturation)
	gm.activeRequests, err = meter.Int64UpDownCounter(
		"grpc.server.active_requests",
		metric.WithDescription("Number of active gRPC requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	return gm, nil
}

// RecordRequest records a completed gRPC request with all golden signals
func (gm *GrpcMetrics) RecordRequest(ctx context.Context, service, method string, duration time.Duration, code codes.Code) {
	attrs := []attribute.KeyValue{
		attribute.String("grpc_service", service),
		attribute.String("grpc_method", method),
		attribute.String("grpc_code", code.String()),
	}

	// Latency
	gm.requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	// Traffic
	gm.requestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Errors (non-OK codes)
	if code != codes.OK {
		gm.errorsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// StartRequest marks the beginning of a request (for saturation tracking)
func (gm *GrpcMetrics) StartRequest(ctx context.Context, service, method string) {
	attrs := []attribute.KeyValue{
		attribute.String("grpc_service", service),
		attribute.String("grpc_method", method),
	}

	gm.activeRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// EndRequest marks the end of a request (for saturation tracking)
func (gm *GrpcMetrics) EndRequest(ctx context.Context, service, method string) {
	if gm == nil || gm.activeRequests == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("grpc_service", service),
		attribute.String("grpc_method", method),
	}

	gm.activeRequests.Add(ctx, -1, metric.WithAttributes(attrs...))
}

// UnaryServerInterceptor returns a gRPC interceptor that records golden signals
func (gm *GrpcMetrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if gm == nil {
			return handler(ctx, req)
		}

		service, method := splitMethodName(info.FullMethod)

		// Track saturation
		gm.StartRequest(ctx, service, method)
		defer gm.EndRequest(ctx, service, method)

		// Record latency
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		// Get gRPC status code
		code := codes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			} else {
				code = codes.Unknown
			}
		}

		gm.RecordRequest(ctx, service, method, duration, code)

		return resp, err
	}
}

// splitMethodName splits "/package.Service/Method" into service and method
func splitMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	if i := strings.LastIndex(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}
	return "unknown", fullMethod
}
