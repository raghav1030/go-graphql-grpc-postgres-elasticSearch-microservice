package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/catalog/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterProductServiceServer(serv, &grpcServer{
		pb.UnimplementedProductServiceServer{},
		s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Price, r.Description, r.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.PostProductResponse{Product: &pb.Product{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	},
	}, nil

}
func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}
func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {

	var res []Product
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Skip, r.Take, r.Query)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	products := []*pb.Product{}

	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
