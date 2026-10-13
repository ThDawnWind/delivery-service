package delivery

import (
	"context"

	deliveryv1 "github.com/ThDawnWind/food-delivery-contracts/gen/go/delivery/v1"
	"google.golang.org/grpc"
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

func (s *GRPCServer) TrackDelivery(
	req *deliveryv1.TrackDeliveryRequest,
	stream grpc.ServerStreamingServer[deliveryv1.DeliveryStatusUpdate],
) error {
	statuses := []deliveryv1.DeliveryStatus{
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_CREATED,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_COURIER_ASSIGNED,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_PICKED_UP,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_ON_THE_WAY,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_DELIVERED,
	}

	for _, deliveryStatus := range statuses {
		err := stream.Send(
			&deliveryv1.DeliveryStatusUpdate{
				DeliveryId: req.GetDeliveryId(),
				Status:     deliveryStatus,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
