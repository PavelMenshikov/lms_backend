package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type EventType string

const (
	NotificationEvent EventType = "events:notifications"
)

type Task struct {
	Type    string            `json:"type"`
	Payload map[string]string `json:"payload"`
}

type EventBroker struct {
	client *redis.Client
}

func NewEventBroker(addr string) *EventBroker {
	return &EventBroker{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (b *EventBroker) Publish(ctx context.Context, event EventType, task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}
	return b.client.LPush(ctx, string(event), data).Err()
}

func (b *EventBroker) Subscribe(ctx context.Context, event EventType) ([]string, error) {
	return b.client.BRPop(ctx, 0, string(event)).Result()
}

func (b *EventBroker) Close() error {
	return b.client.Close()
}
