package main

import (
	"log"

	"github.com/cmrd-a/gophermart/internal/api"
)

func main() {
	r := api.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
