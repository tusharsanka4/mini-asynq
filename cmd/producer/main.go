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

	taskTypes := []string{
		"email:send",
		"image:resize",
		"report:generate",
	}

	// Create and enqueue three tasks.
	for _, taskType := range taskTypes {
		t, err := task.NewTask(
			taskType,
			map[string]string{
				"message": "Hello from mini-asynq!",
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		if err := redisBroker.Enqueue(ctx, t); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Enqueued:", t.Type, t.ID)
	}

	// Retrieve the tasks in FIFO order.
	fmt.Println("\nDequeuing tasks:")

	for range taskTypes {
		t, err := redisBroker.Dequeue(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if t == nil {
			fmt.Println("Queue is empty")
			break
		}

		fmt.Println("Dequeued:", t.Type, t.ID)
	}

}
