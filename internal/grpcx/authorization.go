package grpcx

import (
	"context"
	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func authorizeClient(
	ctx context.Context,
	allowedIdentity string,
) error {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Error(
			codes.Unauthenticated,
			"peer information is missing",
		)
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return status.Error(
			codes.Unauthenticated,
			"TLS authentication information is missing",
		)
	}

	if len(tlsInfo.State.VerifiedChains) == 0 ||
		len(tlsInfo.State.VerifiedChains[0]) == 0 {
		return status.Error(
			codes.Unauthenticated,
			"verified client certificate is missing",
		)
	}

	clientCertificate := tlsInfo.State.VerifiedChains[0][0]

	if !slices.Contains(
		clientCertificate.DNSNames,
		allowedIdentity,
	) {
		return status.Error(
			codes.PermissionDenied,
			"client is not authorized",
		)
	}

	return nil
}

func AuthorizationUnaryInterceptor(
	allowedIdentity string,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if err := authorizeClient(
			ctx,
			allowedIdentity,
		); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func AuthorizationStreamInterceptor(
	allowedIdentity string,
) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if err := authorizeClient(
			stream.Context(),
			allowedIdentity,
		); err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
