package broker

import (
	"context"
	"sync"
)

type MemoryBroker struct {
	subs  map[string][]func([]byte) error
	subMu sync.RWMutex
}

func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{
		subs: make(map[string][]func([]byte) error),
	}
}

func (m *MemoryBroker) Publish(ctx context.Context, topic string, payload []byte) error {
	m.subMu.RLock()
	handlers := m.subs[topic]
	m.subMu.RUnlock()

	for _, h := range handlers {
		go func(fn func([]byte) error) {
			_ = fn(payload)
		}(h)
	}
	return nil
}

func (m *MemoryBroker) Subscribe(topic string, handler func([]byte) error) error {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	m.subs[topic] = append(m.subs[topic], handler)
	return nil
}
