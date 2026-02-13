package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const BoardEventsChannel = "board:events"

type MessageEnvelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Publisher interface {
	Publish(ctx context.Context, msgType string, payload interface{}) error
}

type Broadcaster interface {
	BroadcastRaw(message []byte)
}

type RedisPublisher struct {
	client  *redis.Client
	channel string
}

func NewRedisPublisher(client *redis.Client, channel string) *RedisPublisher {
	return &RedisPublisher{client: client, channel: channel}
}

func (p *RedisPublisher) Publish(ctx context.Context, msgType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("publish marshal payload: %w", err)
	}
	envelope := MessageEnvelope{Type: msgType, Payload: payloadBytes}
	messageBytes, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("publish marshal envelope: %w", err)
	}
	if err := p.client.Publish(ctx, p.channel, messageBytes).Err(); err != nil {
		return fmt.Errorf("publish redis: %w", err)
	}
	return nil
}

type RedisSubscriber struct {
	client  *redis.Client
	hub     Broadcaster
	channel string
}

func NewRedisSubscriber(client *redis.Client, hub Broadcaster, channel string) *RedisSubscriber {
	return &RedisSubscriber{client: client, hub: hub, channel: channel}
}

func (s *RedisSubscriber) Subscribe(ctx context.Context) {
	pubsub := s.client.Subscribe(ctx, s.channel)
	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			_ = pubsub.Close()
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}
			s.hub.BroadcastRaw([]byte(msg.Payload))
		}
	}
}
