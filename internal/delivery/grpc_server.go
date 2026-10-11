package delivery

import (
	"context"

	deliveryv1 "github.com/ThDawnWind/food-delivery-contracts/gen/go/delivery/v1"
)

type GRPCServer struct {
	deliveryv1.UnimplementedDeliveryServiceServer

	service *Service
}

func NewGRPCServer(service *Service) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

func (s *GRPCServer) CalculateDelivery(
	ctx context.Context,
	req *deliveryv1.CalculateDeliveryRequest,
) (*deliveryv1.CalculateDeliveryResponse, error) {
	orderAmount := req.GetOrderAmount()

	deliveryPrice := s.service.Calculate(orderAmount)

	return &deliveryv1.CalculateDeliveryResponse{
		DeliveryPrice: deliveryPrice,
	}, nil
}
