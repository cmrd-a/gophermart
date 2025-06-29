package main

import (
	"context"
	"log"

	"github.com/cmrd-a/gophermart/internal/api"
	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/repository"
	"github.com/cmrd-a/gophermart/internal/service"
)

func main() {
	config.InitConfig()
	repo, _ := repository.NewRepository()
	svc := service.NewService(context.TODO(), *repo)
	r := api.SetupRouter(svc)
	err := r.Run(config.Config.RunAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}
