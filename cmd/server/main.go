package main

import (
	"log"
	"net"

	deliveryv1 "github.com/ThDawnWind/food-delivery-contracts/gen/go/delivery/v1"
	"github.com/ThDawnWind/delivery-service/internal/delivery"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	deliveryService := delivery.NewService()
	deliveryGRPCServer := delivery.NewGRPCServer(deliveryService)
	
	grpcServer := grpc.NewServer()

	deliveryv1.RegisterDeliveryServiceServer(
		grpcServer,
		deliveryGRPCServer,
	)

	log.Println("gRPC server started on :50051")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
