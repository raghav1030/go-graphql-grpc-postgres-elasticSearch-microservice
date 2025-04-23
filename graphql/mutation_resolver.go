// Create Product, Account, Order
package main

import (
	"context"
	"errors"
	"time"

	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/order"
)

var (
	ErrInvalidParameter = errors.New("Invalid Parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	context, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(context, in.Name)
	if err != nil {
		return nil, err
	}
	return &Account{
		a.Id,
		a.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	context, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(context, in.Price, in.Description, in.Name)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	context, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	products := []order.OrderedProduct{}

	for _, op := range in.Products {

		if op.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}

		products = append(products, order.OrderedProduct{
			Id:       op.ID,
			Quantity: uint32(op.Quantity),
		})
	}
	o, err := r.server.orderClient.PostOrder(context, in.AccountID, products)
	if err != nil {
		return nil, err
	}
	return &Order{
		ID:         o.Id,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
	}, nil
}
