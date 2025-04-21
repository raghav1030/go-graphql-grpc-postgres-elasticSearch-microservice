package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/order"
	"github.com/tinrab/kit/retry"
)

type Config struct {
	DatabaseUrl string "envconfig:`DATABASE_URL`"
	CatalogUrl  string "envconfig:`CATALOG_SERVICE_URL`"
	AccountUrl  string "envconfig:`ACCOUNT_SERVICE_URL`"
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = order.NewPostgresRepository(cfg.DatabaseUrl)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	defer r.Close()

	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountUrl, cfg.CatalogUrl, 8080))

}
