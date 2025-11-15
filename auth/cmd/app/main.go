package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	apps "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Println("Shutting down gracefully...")
		cancel()
	}()

	server, err := apps.NewServer(ctx)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	server.Seed()

	server.Run()
}
