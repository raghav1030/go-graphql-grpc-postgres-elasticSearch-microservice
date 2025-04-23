package order

import (
	"context"
	"time"

	pb "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/order/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)
	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderedProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderedProduct{
			ProductId: p.Id,
			Quantity:  p.Quantity,
		})
	}

	res, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountId,
		Products:  protoProducts,
	})

	if err != nil {
		return nil, err
	}

	newOrder := &pb.Order{
		Id:         res.Order.Id,
		AccountId:  accountId,
		TotalPrice: res.Order.TotalPrice,
		Products:   res.Order.Products,
	}
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	return &Order{
		Id:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountId:  accountId,
		Products:   products,
	}, nil

}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	res, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountId,
	})

	if err != nil {
		return nil, err
	}

	orders := []*Order{}

	for _, o := range res.Orders {
		newOrder := &Order{
			Id:         o.Id,
			TotalPrice: o.TotalPrice,
			AccountId:  o.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(o.CreatedAt)

		products := []OrderedProduct{}

		for _, p := range o.Products {
			products = append(products, OrderedProduct{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}

		newOrder.Products = products
		orders = append(orders, newOrder)

	}

	return orders, nil
}
