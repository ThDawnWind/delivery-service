package grpcx

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(
	logger *slog.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		code := status.Code(err)

		requestID := requestIDFromContext(ctx)

		logger.Info(
			"gRPC request completed",
			slog.String("request_id", requestID),
			slog.String("method", info.FullMethod),
			slog.String("code", code.String()),
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

		return resp, err
	}
}

func requestIDFromContext(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(
		ctx,
		"request-id",
	)

	if len(values) == 0 {
		return ""
	}

	return values[0]
}
