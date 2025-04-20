package order

import (
	"context"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, o Order) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type Order struct {
	Id         string           `json:"id"`
	CreatedAt  time.Time        `json:"createdAt"`
	TotalPrice float64          `json:"totalPrice"`
	AccountId  string           `json:"accountId"`
	Products   []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{
		r,
	}
}

func (s *orderService) PostOrder(ctx context.Context, o Order) (*Order, error) {

}
func (s *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {}
