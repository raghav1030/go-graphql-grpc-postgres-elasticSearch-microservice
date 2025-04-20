package main

import (
	"github.com/99designs/gqlgen/graphql"
	account "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/accounts"
)

type Server struct {
		accountClient *account.Client
		orderClient *order.Client
		catalogClient *catalog.Client
}

func NewGraphQLServer(accountUrl string, productUrl string, orderUrl string) (*Server, err) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog.NewClient(productUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil

}

func (s *Server) Mutation() {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) AccountResolver() {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
