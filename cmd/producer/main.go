package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tusharsanka4/mini-asynq/internal/broker"
)

func main() {
	address := os.Getenv("REDIS_ADDR")
	if address == "" {
		address = "localhost:6380"
	}

	fmt.Println("Connecting to Redis...")

	redisBroker := broker.NewBroker(address)
	defer redisBroker.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := redisBroker.Ping(ctx); err != nil {
		log.Fatal("Redis connection failed: ", err)
	}

	fmt.Println("Redis connected successfully!")
}
