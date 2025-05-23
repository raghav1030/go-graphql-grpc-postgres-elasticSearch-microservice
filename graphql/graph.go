package main

import (
	"github.com/99designs/gqlgen/graphql"
	account "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/account"
	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/catalog"
	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/order"
)

type Server struct {
	accountClient account.Client
	orderClient   order.Client
	catalogClient catalog.Client
}

func NewGraphQLServer(accountUrl string, productUrl string, orderUrl string) (*Server, error) {
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
		*accountClient,
		*orderClient,
		*catalogClient,
	}, nil

}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) AccountResolver() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
