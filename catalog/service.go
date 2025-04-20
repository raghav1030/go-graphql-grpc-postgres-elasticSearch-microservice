package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, price float64, description string, name string) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, skip uint64, take uint64, query string) ([]Product, error)
}

type Product struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64    `json:"price"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{
		repository: r,
	}
}

func (s *catalogService) PostProduct(ctx context.Context, price float64, description string, name string) (*Product, error) {
	p := &Product{
		Id:          ksuid.New().String(),
		Name:        name,
		Price:       price,
		Description: description,
	}
	if err := s.repository.PutProduct(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductById(ctx, id)
}
func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.ListProducts(ctx, skip, take)
}
func (s *catalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.ListProductsWithIDs(ctx, ids)
}
func (s *catalogService) SearchProducts(ctx context.Context, skip uint64, take uint64, query string) ([]Product, error) {
	return s.repository.SearchProducts(ctx, skip, take, query)
}
