package main

import (
	"log"
	"time"

	account "github.com/raghav1030/go-graphql-grpc-postgres-elasticSearch-microservice/accounts"
	"github.com/tinrab/kit/retry"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASR_URL"`
}


func main() {
	
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 	var r account.Repository
	// retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
	// 	r, err = account.NewPostgresRepository(cfg.DatabaseURL)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	return
	// })
	// defer r.Close()

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewRepositoryRegistory(cfg.DatabaseUrl)

		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	defer r.Close()

	log.Println("Account Port listening on port 8080...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
