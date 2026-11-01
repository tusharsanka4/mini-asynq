package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tusharsanka4/mini-asynq/internal/task"
)

const defaultQueue = "mini-asynq:queue:default"

type Broker struct {
	Client *redis.Client
}

func NewBroker(address string) *Broker {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return &Broker{
		Client: client,
	}
}

func (b *Broker) Ping(ctx context.Context) error {
	pong, err := b.Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	fmt.Println("PING ->", pong)
	return nil
}

func (b *Broker) Enqueue(ctx context.Context, t *task.Task) error {
	if t == nil || t.ID == "" || t.Type == "" {
		return fmt.Errorf("invalid task")
	}

	// Convert the payload to a string for Redis.
	taskKey := "mini-asynq:task:" + t.ID

	// Atomically store task details and queue its ID.
	_, err := b.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, taskKey, map[string]interface{}{
			"type":       t.Type,
			"payload":    string(t.Payload),
			"created_at": t.CreatedAt.Format(time.RFC3339Nano),
			"status":     "pending",
		})

		pipe.RPush(ctx, defaultQueue, t.ID)
		return nil
	})

	return err
}

func (b *Broker) GetTask(ctx context.Context, id string) (*task.Task, error) {
	taskKey := "mini-asynq:task:" + id

	data, err := b.Client.HGetAll(ctx, taskKey).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("task %s not found", id)
	}

	createdAt, err := time.Parse(time.RFC3339Nano, data["created_at"])
	if err != nil {
		return nil, err
	}

	return &task.Task{
		ID:        id,
		Type:      data["type"],
		Payload:   json.RawMessage(data["payload"]),
		CreatedAt: createdAt,
	}, nil
}

func (b *Broker) Dequeue(ctx context.Context) (*task.Task, error) {
	id, err := b.Client.LPop(ctx, defaultQueue).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return b.GetTask(ctx, id)
}

func (b *Broker) Close() error {
	return b.Client.Close()
}
