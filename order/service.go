package order

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error)
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

func (s *orderService) PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	o := &Order{
		Id:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountId: accountId,
		Products:  products,
	}

	o.TotalPrice = 0.0

	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}

	err := s.repository.PutOrder(ctx, *o)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return o, nil
}
func (s *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountId)
}
