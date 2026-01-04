# Warp Mesh Integration

Go Relay is designed to work seamlessly with **Go Warp**, allowing you to scale from a single in-memory instance to a fully distributed P2P cluster with **zero code changes** in your business logic.

## How it works

The `pkg/adapter/warp` package provides a `Broker` implementation that wraps a Go Warp `MeshNode`.
When you `Enqueue` a job, instead of going to a local channel, it is broadcast to the entire mesh. Any node in the mesh that has subscribed to that topic can pick it up.

## setup

### 1. Install Dependencies
Ensure you have both modules:
```bash
go get github.com/mirkobrombin/go-relay
go get github.com/mirkobrombin/go-warp
```

### 2. Configure the Relay

```go
package main

import (
    "context"
    "log"

    "github.com/mirkobrombin/go-warp/pkg/mesh"
    "github.com/mirkobrombin/go-relay/pkg/manager"
    "github.com/mirkobrombin/go-relay/pkg/adapter/warp"
)

func main() {
    // 1. Start Warp Node (P2P Discovery)
    node, err := mesh.NewNode(mesh.Config{
        ListenAddr: ":0", // Random port
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Wrap it with Relay Adapter
    warpBroker := warp.NewBroker(node)

    // 3. Create Relay Manager with the Warp Broker
    r := manager.New(manager.WithBroker(warpBroker))

    // 4. Register Workers (Same as local!)
    manager.Register(r, "resize_image", ResizeWorker)

    // 5. Start
    r.Start(context.Background())
}
```

## Benefits

*   **Clusterless**: No Redis, RabbitMQ, or Central Broker required.
*   **Decentralized**: Nodes discover each other automatically via Go Warp (UDP/MDNS).
*   **Transparent**: Your `ResizeWorker` function doesn't know it's running on a distributed mesh.
