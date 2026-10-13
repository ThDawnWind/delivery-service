package delivery

import (
	"errors"
	"io"
	"testing"

	deliveryv1 "github.com/ThDawnWind/food-delivery-contracts/gen/go/delivery/v1"
	"google.golang.org/grpc"
)

type fakeTrackDeliveryStream struct {
	grpc.ServerStream
	messages  []*deliveryv1.DeliveryStatusUpdate
	sendError error
}

type fakeCalculateBatchDeliveryStream struct {
	grpc.ServerStream

	requests  []*deliveryv1.CalculateBatchDeliveryRequest
	recvIndex int
	recvError error
	response  *deliveryv1.CalculateBatchDeliveryResponse
}

type fakeTrackDeliveriesStream struct {
	grpc.ServerStream

	requests  []*deliveryv1.TrackDeliveryRequest
	recvIndex int
	responses []*deliveryv1.DeliveryStatusUpdate
	recvError error
	sendError error
}

func (f *fakeCalculateBatchDeliveryStream) Recv() (
	*deliveryv1.CalculateBatchDeliveryRequest,
	error,
) {
	if f.recvError != nil {
		return nil, f.recvError
	}

	if f.recvIndex >= len(f.requests) {
		return nil, io.EOF
	}

	req := f.requests[f.recvIndex]
	f.recvIndex++

	return req, nil
}

func (f *fakeCalculateBatchDeliveryStream) SendAndClose(
	response *deliveryv1.CalculateBatchDeliveryResponse,
) error {
	f.response = response

	return nil
}

func (f *fakeTrackDeliveryStream) Send(
	update *deliveryv1.DeliveryStatusUpdate,
) error {
	if f.sendError != nil {
		return f.sendError
	}

	f.messages = append(f.messages, update)

	return nil
}

func (f *fakeTrackDeliveriesStream) Recv() (
	*deliveryv1.TrackDeliveryRequest,
	error,
) {
	if f.recvError != nil {
		return nil, f.recvError
	}

	if f.recvIndex >= len(f.requests) {
		return nil, io.EOF
	}

	req := f.requests[f.recvIndex]
	f.recvIndex++

	return req, nil
}

func (f *fakeTrackDeliveriesStream) Send(
	response *deliveryv1.DeliveryStatusUpdate,
) error {
	if f.sendError != nil {
		return f.sendError
	}

	f.responses = append(f.responses, response)

	return nil
}

func TestGRPCServer_TrackDelivery(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	stream := &fakeTrackDeliveryStream{}

	req := &deliveryv1.TrackDeliveryRequest{
		DeliveryId: 123,
	}

	err := server.TrackDelivery(req, stream)
	if err != nil {
		t.Fatalf("TrackDelivery() error = %v", err)
	}

	expectedStatuses := []deliveryv1.DeliveryStatus{
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_CREATED,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_COURIER_ASSIGNED,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_PICKED_UP,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_ON_THE_WAY,
		deliveryv1.DeliveryStatus_DELIVERY_STATUS_DELIVERED,
	}

	if len(stream.messages) != len(expectedStatuses) {
		t.Fatalf(
			"messages count = %d, want %d",
			len(stream.messages),
			len(expectedStatuses),
		)
	}

	for i, expectedStatus := range expectedStatuses {
		message := stream.messages[i]

		if message.GetDeliveryId() != 123 {
			t.Errorf(
				"message[%d] delivery_id = %d, want %d",
				i,
				message.GetDeliveryId(),
				123,
			)
		}

		if message.GetStatus() != expectedStatus {
			t.Errorf(
				"message[%d] status = %v, want %v",
				i,
				message.GetStatus(),
				expectedStatus,
			)
		}
	}
}

func TestGRPCServer_TrackDelivery_SendError(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	expectedErr := errors.New("send failed")

	stream := &fakeTrackDeliveryStream{
		sendError: expectedErr,
	}

	req := &deliveryv1.TrackDeliveryRequest{
		DeliveryId: 123,
	}

	err := server.TrackDelivery(req, stream)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"TrackDelivery() error = %v, want %v",
			err,
			expectedErr,
		)
	}

	if len(stream.messages) != 0 {
		t.Fatalf(
			"messages count = %d, want 0",
			len(stream.messages),
		)
	}
}

func TestGRPCServer_CalculateBatchDelivery(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	stream := &fakeCalculateBatchDeliveryStream{
		requests: []*deliveryv1.CalculateBatchDeliveryRequest{
			{
				OrderAmount: 120000,
			},
			{
				OrderAmount: 250000,
			},
			{
				OrderAmount: 80000,
			},
		},
	}

	err := server.CalculateBatchDelivery(stream)
	if err != nil {
		t.Fatalf(
			"CalculateBatchDelivery() error = %v",
			err,
		)
	}

	if stream.response == nil {
		t.Fatal("expected response, got nil")
	}

	if stream.response.GetOrdersCount() != 3 {
		t.Fatalf(
			"orders count = %d, want %d",
			stream.response.GetOrdersCount(),
			3,
		)
	}

	if stream.response.GetTotalDeliveryPrice() != 39800 {
		t.Fatalf(
			"total delivery price = %d, want %d",
			stream.response.GetTotalDeliveryPrice(),
			39800,
		)
	}
}

func TestGRPCServer_CalculateBatchDelivery_RecvError(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	expectedErr := errors.New("receive failed")

	stream := &fakeCalculateBatchDeliveryStream{
		recvError: expectedErr,
	}

	err := server.CalculateBatchDelivery(stream)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"CalculateBatchDelivery() error = %v, want %v",
			err,
			expectedErr,
		)
	}

	if stream.response != nil {
		t.Fatal("expected no response")
	}
}

func TestGRPCServer_TrackDeliveries(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	stream := &fakeTrackDeliveriesStream{
		requests: []*deliveryv1.TrackDeliveryRequest{
			{DeliveryId: 101},
			{DeliveryId: 102},
			{DeliveryId: 103},
		},
	}

	err := server.TrackDeliveries(stream)
	if err != nil {
		t.Fatalf("TrackDeliveries() error = %v", err)
	}

	if len(stream.responses) != 3 {
		t.Fatalf(
			"responses count = %d, want 3",
			len(stream.responses),
		)
	}

	expectedIDs := []int64{101, 102, 103}

	for i, expectedID := range expectedIDs {
		response := stream.responses[i]

		if response.GetDeliveryId() != expectedID {
			t.Errorf(
				"response[%d] delivery_id = %d, want %d",
				i,
				response.GetDeliveryId(),
				expectedID,
			)
		}

		if response.GetStatus() != deliveryv1.DeliveryStatus_DELIVERY_STATUS_CREATED {
			t.Errorf(
				"response[%d] status = %v, want CREATED",
				i,
				response.GetStatus(),
			)
		}
	}
}

func TestGRPCServer_TrackDeliveries_RecvError(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	expectedErr := errors.New("receive failed")

	stream := &fakeTrackDeliveriesStream{
		recvError: expectedErr,
	}

	err := server.TrackDeliveries(stream)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"TrackDeliveries() error = %v, want %v",
			err,
			expectedErr,
		)
	}

	if len(stream.responses) != 0 {
		t.Fatalf(
			"responses count = %d, want 0",
			len(stream.responses),
		)
	}
}

func TestGRPCServer_TrackDeliveries_SendError(t *testing.T) {
	t.Parallel()

	service := NewService()
	server := NewGRPCServer(service)

	expectedErr := errors.New("send failed")

	stream := &fakeTrackDeliveriesStream{
		requests: []*deliveryv1.TrackDeliveryRequest{
			{
				DeliveryId: 101,
			},
		},
		sendError: expectedErr,
	}

	err := server.TrackDeliveries(stream)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"TrackDeliveries() error = %v, want %v",
			err,
			expectedErr,
		)
	}

	if len(stream.responses) != 0 {
		t.Fatalf(
			"responses count = %d, want 0",
			len(stream.responses),
		)
	}
}
