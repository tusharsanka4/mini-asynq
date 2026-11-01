package broker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

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

func (b *Broker) Close() error {
	return b.Client.Close()
}
