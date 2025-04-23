package catalog

import (
	"context"
	"log"

	pb "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, price float64, description string, name string) (*Product, error) {
	res, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Product{
		Id:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	res, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Product{
		Id:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}
func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {
	res, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	products := []Product{}

	for _, res := range res.Products {
		products = append(products, Product{
			Id:          res.Id,
			Description: res.Description,
			Name:        res.Name,
			Price:       res.Price,
		})
	}
	return products, nil

}
