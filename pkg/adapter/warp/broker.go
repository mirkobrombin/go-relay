package warp

import (
	"context"
	"fmt"
)

// MeshProvider defines the interface required from a go-warp node.
type MeshProvider interface {
	Broadcast(topic string, data []byte) error
	Listen(topic string, handler func([]byte)) error
}

// Broker adapts a Warp Mesh provider to the Relay Broker interface.
type Broker struct {
	mesh MeshProvider
}

// NewBroker creates a Relay broker backed by a Warp Mesh provider.
func NewBroker(provider MeshProvider) *Broker {
	return &Broker{mesh: provider}
}

func (b *Broker) Publish(ctx context.Context, topic string, payload []byte) error {
	return b.mesh.Broadcast(topic, payload)
}

func (b *Broker) Subscribe(topic string, handler func([]byte) error) error {
	return b.mesh.Listen(topic, func(data []byte) {
		if err := handler(data); err != nil {
			fmt.Printf("[WarpBroker] Handler execution failed: %v\n", err)
		}
	})
}
