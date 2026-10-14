package grpcx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func TestAuthorizeClient_Allowed(t *testing.T) {
	t.Parallel()

	certificate := &x509.Certificate{
		DNSNames: []string{
			"food-delivery-api",
		},
	}

	ctx := peer.NewContext(
		context.Background(),
		&peer.Peer{
			AuthInfo: credentials.TLSInfo{
				State: tls.ConnectionState{
					VerifiedChains: [][]*x509.Certificate{
						{
							certificate,
						},
					},
				},
			},
		},
	)

	err := authorizeClient(
		ctx,
		"food-delivery-api",
	)
	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}
}

func TestAuthorizeClient_PermissionDenied(t *testing.T) {
	t.Parallel()

	certificate := &x509.Certificate{
		DNSNames: []string{
			"analytics-service",
		},
	}

	ctx := peer.NewContext(
		context.Background(),
		&peer.Peer{
			AuthInfo: credentials.TLSInfo{
				State: tls.ConnectionState{
					VerifiedChains: [][]*x509.Certificate{
						{
							certificate,
						},
					},
				},
			},
		},
	)

	err := authorizeClient(
		ctx,
		"food-delivery-api",
	)
	if err == nil {
		t.Fatal(
			"expected permission denied error, got nil",
		)
	}

	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf(
			"expected code %s, got %s",
			codes.PermissionDenied,
			status.Code(err),
		)
	}
}

func TestAuthorizeClient_MissingPeer(t *testing.T) {
	t.Parallel()

	err := authorizeClient(
		context.Background(),
		"food-delivery-api",
	)
	if err == nil {
		t.Fatal(
			"expected unauthenticated error, got nil",
		)
	}

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf(
			"expected code %s, got %s",
			codes.Unauthenticated,
			status.Code(err),
		)
	}
}

func TestAuthorizeClient_MissingTLSInfo(t *testing.T) {
	t.Parallel()

	ctx := peer.NewContext(
		context.Background(),
		&peer.Peer{},
	)

	err := authorizeClient(
		ctx,
		"food-delivery-api",
	)
	if err == nil {
		t.Fatal(
			"expected unauthenticated error, got nil",
		)
	}

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf(
			"expected code %s, got %s",
			codes.Unauthenticated,
			status.Code(err),
		)
	}
}
