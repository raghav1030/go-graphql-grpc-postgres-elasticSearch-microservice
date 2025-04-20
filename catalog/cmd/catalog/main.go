package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/catalog"
	"github.com/tinrab/kit/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}
	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseUrl)
		if err != nil {
			log.Fatal(err)
		}
		return
	})
	r.Close()

	log.Println("Catalog Service listening on port 8080...")

	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
