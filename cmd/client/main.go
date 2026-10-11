package main

import (
	"context"
	"log"
	"time"

	deliveryv1 "github.com/ThDawnWind/food-delivery-contracts/gen/go/delivery/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := deliveryv1.NewDeliveryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.CalculateDelivery(
		ctx,
		&deliveryv1.CalculateDeliveryRequest{
			OrderAmount: 250000,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("delivery price: %d", resp.GetDeliveryPrice())
}
