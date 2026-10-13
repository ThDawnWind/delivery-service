package main

import (
	"log"
	"log/slog"
	"net"

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

	tlsCredentials, err := credentials.NewServerTLSFromFile(
		"certs/server.crt",
		"certs/server.key",
	)
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(
			grpcx.LoggingInterceptor(slog.Default()),
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
