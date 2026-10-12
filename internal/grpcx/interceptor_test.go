package grpcx

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestLoggingInterceptor(t *testing.T) {
	t.Parallel()

	logger := slog.New(
		slog.NewTextHandler(
			io.Discard,
			nil,
		),
	)

	interceptor := LoggingInterceptor(logger)

	handlerCalled := false

	handler := func(
		_ context.Context,
		req any,
	) (any, error) {
		handlerCalled = true

		return req, nil
	}

	req := "test-request"

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	resp, err := interceptor(
		context.Background(),
		req,
		info,
		handler,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !handlerCalled {
		t.Fatal("expected handler to be called")
	}

	if resp != req {
		t.Fatalf(
			"response = %v, want %v",
			resp,
			req,
		)
	}
}

func TestLoggingInterceptor_ReturnsHandlerError(t *testing.T) {
	t.Parallel()

	logger := slog.New(
		slog.NewTextHandler(
			io.Discard,
			nil,
		),
	)

	interceptor := LoggingInterceptor(logger)

	handler := func(
		_ context.Context,
		_ any,
	) (any, error) {
		return nil, status.Error(
			codes.Unavailable,
			"service unavailable",
		)
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	resp, err := interceptor(
		context.Background(),
		"test-request",
		info,
		handler,
	)

	if resp != nil {
		t.Fatalf(
			"response = %v, want nil",
			resp,
		)
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if status.Code(err) != codes.Unavailable {
		t.Fatalf(
			"status code = %s, want %s",
			status.Code(err),
			codes.Unavailable,
		)
	}
}

func TestRequestIDFromContext(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			"request-id",
			"test-request-id",
		),
	)

	got := requestIDFromContext(ctx)

	if got != "test-request-id" {
		t.Fatalf(
			"request ID = %q, want %q",
			got,
			"test-request-id",
		)
	}
}

func TestRequestIDFromContext_NoMetadata(t *testing.T) {
	t.Parallel()

	got := requestIDFromContext(
		context.Background(),
	)

	if got != "" {
		t.Fatalf(
			"request ID = %q, want empty string",
			got,
		)
	}
}
