package usecase

import (
	"context"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute(ctx context.Context) ([]OrderOutputDTO, error) {
	orders, err := c.OrderRepository.List(ctx)
	if err != nil {
		return nil, err
	}
	payloadOrders := []OrderOutputDTO{}
	for _, order := range orders {
		payloadOrders = append(payloadOrders, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}
	return payloadOrders, nil
}
