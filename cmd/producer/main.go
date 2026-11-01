package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tusharsanka4/mini-asynq/internal/broker"
	"github.com/tusharsanka4/mini-asynq/internal/task"
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

	emailTask, err := task.NewTask(
		"email:send",
		map[string]string{
			"to":      "bob@example.com",
			"subject": "Hello from mini-asynq",
			"body":    "This is our first task!",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	encoded, err := emailTask.Serialize()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSerialized task:")
	fmt.Println(string(encoded))

	decoded, err := task.Deserialize(encoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDeserialized task:")
	fmt.Println("ID:", decoded.ID)
	fmt.Println("Type:", decoded.Type)
	fmt.Println("Payload:", string(decoded.Payload))
	fmt.Println("Created At:", decoded.CreatedAt)

}
