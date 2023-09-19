package service

import (
	"context"
	"time"

	"github.com/nimbo1999/20-CleanArch/internal/infra/grpc/pb"
	"github.com/nimbo1999/20-CleanArch/internal/usecase"
)

// type OrderServiceClient interface {
// 	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
// 	ListOrders(ctx context.Context, in *Null, opts ...grpc.CallOption) (OrderService_ListOrdersClient, error)
// }

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(in *pb.Null, stream pb.OrderService_ListOrdersServer) error {
	ctx, cancel := context.WithDeadline(stream.Context(), time.Now().Add(time.Second))
	defer cancel()
	orders, err := s.ListOrderUseCase.Execute(ctx)

	if err != nil {
		return err
	}

	for _, order := range orders {
		err = stream.Send(&pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
