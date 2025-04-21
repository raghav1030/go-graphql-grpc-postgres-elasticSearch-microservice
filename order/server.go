package order

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/order/pb"

	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/account"
	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	accountClient account.Client
	catalogClient catalog.Client
	service       Service
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, catalogUrl string, accountUrl string, port uint32) error {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))

	if err != nil {
		accountClient.Close()
		catalogClient.Close()
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		*accountClient,
		*catalogClient,
		s,
		pb.UnimplementedOrderServiceServer{},
	})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Fatal("Error getting account", nil)
		return nil, err
	}

	productIDs := []string{}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Fatal("Error fetching ordered products", err)
		return nil, err
	}

	products := []OrderedProduct{}

	for _, p := range orderedProducts {
		product := &OrderedProduct{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.Id {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, *product)
		}

	}
	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         order.Id,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderedProduct{},
	}

	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderedProduct{
			Id:          p.Id,
			Description: p.Description,
			Name:        p.Name,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})

	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

// func (s *grpcServer) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {

// }

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Fatal("Error fetching account orders")
		return nil, err
	}
	productIdMap := map[string]bool{}

	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIdMap[p.Id] = true
		}
	}

	producIDs := []string{}

	for id := range productIdMap {
		producIDs = append(producIDs, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, producIDs, "")
	if err != nil {
		log.Fatal("Cannot fetch products from product IDs")
		return nil, err
	}
	orders := []*pb.Order{}
	for _, o := range accountOrders {
		op := &pb.Order{
			Id:         o.Id,
			AccountId:  o.AccountId,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderedProduct{},
		}
		op.CreatedAt, err = o.CreatedAt.MarshalBinary()

		if err != nil {
			log.Fatal("Error converting created at into binary object")
			return nil, err
		}
		for _, product := range o.Products {
			for _, p := range products {
				if p.Id == product.Id {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}
			op.Products = append(op.Products, &pb.Order_OrderedProduct{
				Id:          product.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}
		orders = append(orders, op)

	}

	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
