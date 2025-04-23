// Having order Query because Order is dependent on Account id so it is not a root resolver
package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, acc *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, acc.ID)

	if err != nil {
		return nil, err
	}

	orders := []*Order{}

	for _, o := range orderList {
		products := []*OrderedProduct{}

		for _, p := range o.Products {

			products = append(products, &OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}

		orders = append(orders, &Order{
			ID:         o.Id,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
		})
	}
	return orders, nil
}
