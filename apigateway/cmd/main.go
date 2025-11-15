package main

import (
	"log"
	"os"
	"os/signal"

	apps "github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/app"
)

func main() {
	_, shutdown, err := apps.RunClient()
	if err != nil {
		log.Fatalf("failed to run client: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	shutdown()
}
