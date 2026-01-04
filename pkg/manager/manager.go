package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mirkobrombin/go-relay/pkg/broker"
)

type Job struct {
	ID        string
	Queue     string
	Topic     string
	Payload   []byte
	CreatedAt time.Time
	TryCount  int
}

type Handler[T any] func(ctx context.Context, payload T) error

type Broker interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(topic string, handler func(payload []byte) error) error
}

type Relay struct {
	broker    Broker
	handlers  map[string]any
	handlerMu sync.RWMutex
}

type Option func(*Relay)

func New(opts ...Option) *Relay {
	r := &Relay{
		broker:   broker.NewMemoryBroker(),
		handlers: make(map[string]any),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithBroker(b Broker) Option {
	return func(r *Relay) {
		r.broker = b
	}
}

func Register[T any](r *Relay, topic string, fn Handler[T]) {
	wrapper := func(raw []byte) error {
		var payload T
		if err := json.Unmarshal(raw, &payload); err != nil {
			return fmt.Errorf("payload unmarshal failed: %w", err)
		}
		return fn(context.Background(), payload)
	}

	r.handlerMu.Lock()
	r.handlers[topic] = wrapper
	r.handlerMu.Unlock()
}

func Enqueue[T any](ctx context.Context, r *Relay, topic string, payload T) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload marshal failed: %w", err)
	}

	return r.broker.Publish(ctx, topic, data)
}

func (r *Relay) Start(ctx context.Context) error {
	r.handlerMu.RLock()
	defer r.handlerMu.RUnlock()

	for topic, wrapperFn := range r.handlers {
		userHandler := wrapperFn.(func([]byte) error)

		err := r.broker.Subscribe(topic, func(data []byte) error {
			defer func() {
				if rec := recover(); rec != nil {
					fmt.Printf("panic in job %s: %v\n", topic, rec)
				}
			}()

			return userHandler(data)
		})
		if err != nil {
			return err
		}
	}

	<-ctx.Done()
	return nil
}
