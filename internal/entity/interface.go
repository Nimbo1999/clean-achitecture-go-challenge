package entity

import "context"

type OrderRepositoryInterface interface {
	Save(order *Order) error
	List(ctx context.Context) ([]*Order, error)
}
