package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Broker struct {
	client *redis.Client
}

func NewBroker(client *redis.Client) *Broker {
	return &Broker{client: client}
}

func (b *Broker) PublishJSON(ctx context.Context, channel string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, channel, data).Err()
}

func (b *Broker) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return b.client.Subscribe(ctx, channel)
}

func (b *Broker) IncrementCounter(ctx context.Context, key string, delta int64) (int64, error) {
	return b.client.IncrBy(ctx, key, delta).Result()
}

func (b *Broker) SetCounter(ctx context.Context, key string, value int64, ttl time.Duration) error {
	return b.client.Set(ctx, key, value, ttl).Err()
}

func (b *Broker) GetCounter(ctx context.Context, key string) (int64, error) {
	return b.client.Get(ctx, key).Int64()
}

func (b *Broker) ResetCounter(ctx context.Context, key string) error {
	return b.client.Set(ctx, key, 0, 24*time.Hour).Err()
}

func ChatChannel(orderID string) string {
	return fmt.Sprintf("chat:order:%s", orderID)
}

func NotificationsChannel(userID string) string {
	return fmt.Sprintf("notifications:user:%s", userID)
}

func ChatUnreadKey(userID, orderID string) string {
	return fmt.Sprintf("chat:unread:%s:%s", userID, orderID)
}

func NotificationsUnreadKey(userID string) string {
	return fmt.Sprintf("notifications:unread:%s", userID)
}
