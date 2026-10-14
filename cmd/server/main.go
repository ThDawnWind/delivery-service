package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/ThDawnWind/delivery-service/internal/delivery"
	"github.com/ThDawnWind/delivery-service/internal/grpcx"
	deliveryv1 "github.com/ThDawnWind/food-delivery-contracts/gen/go/delivery/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	deliveryService := delivery.NewService()
	deliveryGRPCServer := delivery.NewGRPCServer(deliveryService)

	serverCertificate, err := tls.LoadX509KeyPair(
		"certs/server.crt",
		"certs/server.key",
	)
	if err != nil {
		log.Fatalf("failed to load server TLS certificate: %v", err)
	}

	caCertificate, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}

	clientCAPool := x509.NewCertPool()

	if ok := clientCAPool.AppendCertsFromPEM(caCertificate); !ok {
		log.Fatal("failed to add CA certificate to client CA pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			serverCertificate,
		},
		ClientCAs:  clientCAPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	tlsCredentials := credentials.NewTLS(tlsConfig)

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),

		grpc.ChainUnaryInterceptor(
			grpcx.LoggingInterceptor(slog.Default()),
			grpcx.AuthorizationUnaryInterceptor(
				"food-delivery-api",
			),
		),

		grpc.ChainStreamInterceptor(
			grpcx.AuthorizationStreamInterceptor(
				"food-delivery-api",
			),
		),
	)

	deliveryv1.RegisterDeliveryServiceServer(
		grpcServer,
		deliveryGRPCServer,
	)

	reflection.Register(grpcServer)

	healthServer := health.NewServer()

	healthServer.SetServingStatus(
		"delivery.v1.DeliveryService",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

	grpc_health_v1.RegisterHealthServer(
		grpcServer,
		healthServer,
	)

	log.Println("gRPC server started on :50051")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
